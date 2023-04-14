package main

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	_ "net/http/pprof"
	"selfadaptive/controllers/def/info"
	"selfadaptive/rabbitmq/adaptationlogic"
	"selfadaptive/rabbitmq/mytimer"
	"selfadaptive/shared"
	"time"
)

type Subscriber struct {
	IsAdaptive bool
	Conn       *amqp.Connection
	Ch         *amqp.Channel
	Queue      amqp.Queue
	Msgs       <-chan amqp.Delivery
	PC         int
}

func main() {

	// configure/read flags
	var isAdaptivePtr = flag.Bool("is-adaptive", false, "is-adaptive is a boolean")
	var controllerTypePtr = flag.String("controller-type", "OnOff", "controller-type is a string")
	var monitorIntervalPtr = flag.Int("monitor-interval", 1, "monitor-interval is an int (s)")
	var setPointPtr = flag.Float64("set-point", 3000.0, "set-point is a float (goal rate)")
	var kpPtr = flag.Float64("kp", 1.0, "Kp is a float")
	var kiPtr = flag.Float64("ki", 1.0, "Ki is a float")
	var kdPtr = flag.Float64("kd", 1.0, "Kd is a float")
	var prefetchCountPtr = flag.Int("prefetch-count", 1, "prefetch-count is an int")
	var minPtr = flag.Float64("min", 0.0, "min is a float")
	var maxPtr = flag.Float64("max", 100.0, "max is a float")
	var deadZonePtr = flag.Float64("dead-zone", 0.0, "dead-zone is a float")
	var hysteresisBandPtr = flag.Float64("hysteresis-band", 0.0, "hysteresis-band is a float")
	var directionPtr = flag.Float64("direction", 1.0, "direction is a float")
	flag.Parse()

	// create new consumer
	var consumer = NewConsumer(*isAdaptivePtr, *prefetchCountPtr)

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

	fmt.Println("*************** Subscriber started *************")
	fmt.Printf("IsAdaptive     = % v\n", *isAdaptivePtr)
	fmt.Printf("Controller Type= %v\n", *controllerTypePtr)
	fmt.Printf("Direction      = %.1f\n", *directionPtr)
	fmt.Printf("Min            = %.2f\n", *minPtr)
	fmt.Printf("Max            = %.2f\n", *maxPtr)
	fmt.Printf("Kp             = %.2f\n", *kpPtr)
	fmt.Printf("Ki             = %.2f\n", *kiPtr)
	fmt.Printf("Kd             = %.2f\n", *kdPtr)
	fmt.Printf("Dead Zone      = %.2f\n", *deadZonePtr)
	fmt.Printf("Hysteresis Band= %.2f\n", *hysteresisBandPtr)
	fmt.Printf("Goal           = %.2f\n", *setPointPtr)
	fmt.Printf("Monitor Time   = %v (s)\n", *monitorIntervalPtr)
	fmt.Printf("Prefetch Count = %v\n", *prefetchCountPtr)
	fmt.Println("************************************************")

	if *isAdaptivePtr {

		// Create & start adaptation logic
		c := info.Controller{TypeName: *controllerTypePtr, Direction: *directionPtr, Min: *minPtr, Max: *maxPtr, Kp: *kpPtr, Ki: *kiPtr, Kd: *kdPtr, DeadZone: *deadZonePtr, HysteresisBand: *hysteresisBandPtr}
		adapter := adaptationlogic.NewAdaptationLogic(toAdapter, fromAdapter, c, *setPointPtr, time.Duration(*monitorIntervalPtr), *prefetchCountPtr)
		go adapter.Run()

		// Create timer
		t := mytimer.NewMyTimer(*monitorIntervalPtr, startTimer, stopTimer)
		go t.RunMyTimer()

		// run adaptive consumer
		consumer.RunAdaptive(startTimer, stopTimer, toAdapter, fromAdapter)
	} else {
		// run non-adaptive consumer
		consumer.RunNonAdaptive()
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

func (c Subscriber) RunAdaptive(startTimer, stopTimer chan bool, toAdapter chan int, fromAdapter chan int) {

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
	c.Conn, err = amqp.Dial("amqp://guest:guest@192.168.0.20:5672/") // Home Recife
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
