package main

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"main.go/controllers/def/info"
	"main.go/rabbitmq/adaptationlogic"
	"main.go/rabbitmq/mytimer"
	"main.go/shared"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"
)

type ExecutionParameters struct {
	Tunning         *string
	ExecutionType   *string
	IsAdaptive      *bool
	ControllerType  *string
	MonitorInterval *int
	SetPoint        *float64
	Kp              *float64
	Ki              *float64
	Kd              *float64
	PrefetchCount   *int
	Min             *float64
	Max             *float64
	DeadZone        *float64
	HysteresisBand  *float64
	Direction       *float64
	GainTrigger     *float64
	Beta            *float64
}

type Subscriber struct {
	IsAdaptive bool
	Conn       *amqp.Connection
	Ch         *amqp.Channel
	Queue      amqp.Queue
	Msgs       <-chan amqp.Delivery
	PC         int
}

func main() {

	runtime.GOMAXPROCS(1) // TODO

	// load parameters
	p := loadParameters()

	// validate parameters
	validateParameters(p) // TODO

	// show parameters
	showParameters(p)

	// create new consumer
	var consumer = NewConsumer(*p.IsAdaptive, *p.PrefetchCount)

	// Configure RabbitMQ
	consumer.ConfigureRabbitMQ(consumer.PC)
	defer func(Conn *amqp.Connection) {
		err := Conn.Close()
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}(consumer.Conn)
	defer func(Ch *amqp.Channel) {
		err := Ch.Close()
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}(consumer.Ch)

	// create channels
	toAdapter := make(chan int)   // no. of messages
	fromAdapter := make(chan int) // pc
	startTimer := make(chan bool) // start timer
	stopTimer := make(chan bool)  // stop timer

	if *p.IsAdaptive {

		// Create & start adaptation logic
		c := info.Controller{TypeName: *p.ControllerType, Direction: *p.Direction, PC: float64(*p.PrefetchCount), Min: *p.Min, Max: *p.Max, Kp: *p.Kp, Ki: *p.Ki, Kd: *p.Kd, DeadZone: *p.DeadZone, HysteresisBand: *p.HysteresisBand, GainTrigger: *p.GainTrigger, Beta: *p.Beta}
		adapter := adaptationlogic.NewAdaptationLogic(*p.ExecutionType, toAdapter, fromAdapter, c, *p.SetPoint, time.Duration(*p.MonitorInterval), *p.PrefetchCount)
		go adapter.Run() // normal execution

		// Create timer
		t := mytimer.NewMyTimer(*p.MonitorInterval, startTimer, stopTimer)
		go t.RunMyTimer()

		// run adaptive consumer
		consumer.RunAdaptive(startTimer, stopTimer, toAdapter, fromAdapter)
	} else {
		// run non-adaptive consumer
		//consumer.RunNonAdaptive()
		// Create timer
		t := mytimer.NewMyTimer(*p.MonitorInterval, startTimer, stopTimer)
		go t.RunMyTimer()
		consumer.RunNonAdaptiveMonitored(startTimer, stopTimer, p)
	}
}

func (c Subscriber) RunNonAdaptive() {

	for {
		d := <-c.Msgs
		err := d.Ack(false) // send ack to broker
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}
}

func (c Subscriber) RunNonAdaptiveMonitored(startTimer, stopTimer chan bool, p ExecutionParameters) {
	count := 0
	for d := range c.Msgs {
		err := d.Ack(false) // send ack to broker
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		count++ // increment number of received messages
		select {
		case <-startTimer: // start monitor timer
			// initialise no. of receive messages
			count = 0
		case <-stopTimer: // stop monitor timer
			rate := count / *p.MonitorInterval
			fmt.Println(0, ";", *p.PrefetchCount, ";", rate)
			// inspect queue
			/*q, err1 := c.Ch.QueueInspect("rpc_queue")
			if err1 != nil {
				shared.ErrorHandler(shared.GetFunction(), "Impossible to inspect the queue")
				os.Exit(0)
			}
			*/
		default:
		}
	}
}

func (c Subscriber) RunAdaptive(startTimer, stopTimer chan bool, toAdapter chan int, fromAdapter chan int) {

	count := 0

	for d := range c.Msgs {
		err := d.Ack(false) // send ack to broker
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		count++ // increment number of received messages
		select {
		case <-startTimer: // start monitor timer
			// initialise no. of receive messages
			count = 0
		case <-stopTimer: // stop monitor timer
			// send no. of messages to adaptation logic
			toAdapter <- count

			// receive new pc from adaptation logic
			c.PC = <-fromAdapter

			// inspect queue
			/*q, err1 := c.Ch.QueueInspect("rpc_queue")
			if err1 != nil {
				shared.ErrorHandler(shared.GetFunction(), "Impossible to inspect the queue")
				os.Exit(0)
			}
			*/

			// configure new pc
			err := c.Ch.Qos(
				c.PC, // prefetch count
				0,    // prefetch size
				true, // global TODO default is false
			)
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), "Failed to set QoS")
			}
		default:
		}
	}
}

func (c Subscriber) RunAdaptiveOld(startTimer, stopTimer chan bool, toAdapter chan int, fromAdapter chan int) {

	count := 0

	for {
		select {
		case <-startTimer: // start monitor timer
			// initialise no. of receive messages
			count = 0
		case <-stopTimer: // stop monitor timer
			// send no. of messages to adaptation logic
			toAdapter <- count

			// receive new pc from adaptation logic
			c.PC = <-fromAdapter

			// inspect queue
			/*q, err1 := c.Ch.QueueInspect("rpc_queue")
			if err1 != nil {
				shared.ErrorHandler(shared.GetFunction(), "Impossible to inspect the queue")
				os.Exit(0)
			}
			*/

			// configure new pc
			err := c.Ch.Qos(
				c.PC, // prefetch count
				0,    // prefetch size
				true, // global TODO default is false
			)
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), "Failed to set QoS")
			}
		default: // receive a message
			d := <-c.Msgs       // receive a message
			err := d.Ack(false) // send ack to broker
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}
			count++ // increment number of received messages
		}
	}
}

func (c *Subscriber) ConfigureRabbitMQ(pc int) {
	err := error(nil)

	// create connection
	//c.Conn, err = amqp.Dial("amqp://guest:guest@10.45.21.246:5672/") // KU Leuven
	c.Conn, err = amqp.Dial("amqp://guest:guest@" + shared.IpPortRabbitMQ + "/")
	//c.Conn, err = amqp.Dial("amqp://guest:guest@192.168.1.127:5672/") // Leuven
	//c.Conn, err = amqp.Dial("amqp://guest:guest@192.168.0.110:5672/") // Home
	//s.Conn, err = amqp.Dial("amqp://guest:guest@172.22.38.75:5672/") // Home

	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to connect to RabbitMQ")
	}

	//connSub, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	//s.ConnSub, err = amqp.Dial("amqp://nsr:nsr@localhost:5672/") // Docker 'some-rabbit'
	//shared.FailOnError(err, "Failed to connect to RabbitMQ - Subscriber")
	//defer conn.Close()

	// create channel
	c.Ch, err = c.Conn.Channel()
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to open a channel")
	}

	// declare queues
	c.Queue, err = c.Ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to declare a queue")
	}

	// create a consumer
	c.Msgs, err = c.Ch.Consume(
		c.Queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to register a consumer")
	}

	// configure initial QoS of Req channel
	err = c.Ch.Qos(
		pc,   // prefetch count
		0,    // prefetch size
		true, // global TODO default is false
	)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to set QoS")
	}
	return
}

func NewConsumer(isAdaptive bool, pc int) Subscriber {
	s := Subscriber{IsAdaptive: isAdaptive, PC: pc}

	return s
}

func loadParameters() ExecutionParameters {
	p := ExecutionParameters{}

	p.ExecutionType = flag.String("execution-type", shared.StaticGoal, "execution-type is a string")
	p.IsAdaptive = flag.Bool("is-adaptive", false, "is-adaptive is a boolean")
	p.ControllerType = flag.String("controller-type", "OnOff", "controller-type is a string")
	p.MonitorInterval = flag.Int("monitor-interval", 1, "monitor-interval is an int (s)")
	p.SetPoint = flag.Float64("set-point", 3000.0, "set-point is a float (goal rate)")
	p.Kp = flag.Float64("kp", 1.0, "Kp is a float")
	p.Ki = flag.Float64("ki", 1.0, "Ki is a float")
	p.Kd = flag.Float64("kd", 1.0, "Kd is a float")
	p.PrefetchCount = flag.Int("prefetch-count", 1, "prefetch-count is an int")
	p.Min = flag.Float64("min", 0.0, "min is a float")
	p.Max = flag.Float64("max", 100.0, "max is a float")
	p.DeadZone = flag.Float64("dead-zone", 0.0, "dead-zone is a float")
	p.HysteresisBand = flag.Float64("hysteresis-band", 0.0, "hysteresis-band is a float")
	p.Direction = flag.Float64("direction", 1.0, "direction is a float")
	p.GainTrigger = flag.Float64("gain-trigger", 1.0, "gain trigger is a float")
	p.Beta = flag.Float64("beta", 1.0, "Beta is a float (used in PI controllers with two degrees of freedom")
	p.Tunning = flag.String("tunning", "RootLocus", "tunning-type is a string")
	flag.Parse()

	return p
}

func validateParameters(p ExecutionParameters) {
	if *p.Direction != 1.0 && *p.Direction != -1.0 {
		shared.ErrorHandler(shared.GetFunction(), "Direction invalid")
	}

	if *p.Tunning != shared.RootLocus && *p.Tunning != shared.Ziegler && *p.Tunning != shared.Cohen && *p.Tunning != shared.Amigo {
		shared.ErrorHandler(shared.GetFunction(), "Tunning ´"+*p.Tunning+"´ is invalid")
	}
}

func showParameters(p ExecutionParameters) {

	// validate execution type
	fmt.Println("************************************************")
	fmt.Printf("Execution Type  : %v\n", *p.ExecutionType)
	fmt.Printf("Is Adaptive?    : %v\n", *p.IsAdaptive)
	fmt.Printf("Tunning         : %v\n", *p.Tunning)
	fmt.Printf("Controller Type : %v\n", *p.ControllerType)
	fmt.Printf("Monitor Interval: %v\n", *p.MonitorInterval)
	fmt.Printf("Goal            : %.4f\n", *p.SetPoint)
	fmt.Printf("Prefetch Count  : %v\n", *p.PrefetchCount)
	fmt.Printf("Direction       : %.1f\n", *p.Direction)

	switch *p.ControllerType {
	case shared.AsTAR:
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
		fmt.Printf("Hysteresis Band : %.4f\n", *p.HysteresisBand)
	case shared.BasicOnoff:
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
	case shared.DeadZoneOnoff:
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
		fmt.Printf("Dead Zone       : %.4f\n", *p.DeadZone)
	case shared.HysteresisOnoff:
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
		fmt.Printf("Hystereis Band  : %.4f\n", *p.HysteresisBand)
	case shared.BasicP:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
	case shared.BasicPi:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
	case shared.BasicPid:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
	case shared.SmoothingPid:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
	case shared.IncrementalFormPid:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
	case shared.ErrorSquarePid:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
	case shared.DeadZonePid:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
		fmt.Printf("Dead Zone       : %.4f\n", *p.DeadZone)
	case shared.GainScheduling:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
		fmt.Printf("Gain Trigger    : %.4f\n", *p.GainTrigger)
	case shared.HPA:
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
		fmt.Printf("PC           : %v\n", *p.PrefetchCount)
	case shared.PIwithTwoDegreesOfFreedom:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
		fmt.Printf("Beta            : %.4f\n", *p.Beta)
	case shared.WindUp:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
	default:
		fmt.Println(shared.GetFunction(), "Controller type ´", *p.ControllerType, "´ is invalid")
		os.Exit(0)
	}
	fmt.Println("************************************************")
}
