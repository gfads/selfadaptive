package main

/*
import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	_ "net/http/pprof"
	"rabbitmq/journal/actuator"
	"rabbitmq/journal/adjustmentmechanism"
	"rabbitmq/journal/controller"
	"rabbitmq/journal/monitor"
	"rabbitmq/shared"
	"sync/atomic"
	"time"
)

type Consumer struct {
	IsAdaptive bool
	Conn       *amqp.Connection
	Ch         *amqp.Channel
	Queue      amqp.Queue
	Msgs       <-chan amqp.Delivery
}

var count uint64

func oldMain() {

	// configure/read flags
	var isAdaptivePtr = flag.Bool("is-adaptive", false, "is-adaptive is a boolean")
	var controllerTypePtr = flag.String("controller-type", "OnOff", "controller-type is a string")
	var monitorIntervalPtr = flag.Int("monitor-interval", 1, "monitor-interval is an int (s)")
	var setPointPtr = flag.Float64("set-point", 3000.0, "set-point is a float (goal rate)")
	var kpPtr = flag.Float64("kp", 1.0, "Kp is a float")
	var kiPtr = flag.Float64("ki", 1.0, "Ki is a float")
	var kdPtr = flag.Float64("kd", 1.0, "Kd is a float")
	var prefetchCountPtr = flag.Int("prefetch-count", 1, "prefetch-count is an int")
	flag.Parse()

	// create controller
	var c controller.IController
	c = controller.NewController(*controllerTypePtr, *kpPtr, *kiPtr, *kdPtr)

	// Create channel Monitor - Actuator
	ch := make(chan bool)

	// create new consumer
	var consumer = NewConsumer(*isAdaptivePtr, *monitorIntervalPtr, c, *prefetchCountPtr, *setPointPtr)

	// Configure RabbitMQ
	consumer.configureRabbitMQ(*prefetchCountPtr)
	defer consumer.Conn.Close()
	defer consumer.Ch.Close()

	// create & start Monitor
	m := monitor.NewMonitor(time.Duration(*monitorIntervalPtr) * time.Second)
	go m.Run(ch)

	// create & start adjustment mechanism
	startOuter := make(chan bool)
	stopOuter := make(chan bool)
	am := adjustmentmechanism.NewAdjustmentMechanism(time.Duration(*monitorIntervalPtr) * time.Second * 2)
	go am.Run(c, startOuter, stopOuter)

	// create Actuator && configure channels
	startInner := make(chan bool)
	stopInner := make(chan bool)
	a := actuator.NewActuator(*prefetchCountPtr)
	if *isAdaptivePtr {
		//go a.RunStatic(*monitorIntervalPtr, ch, *consumer.Ch, c, *setPointPtr, &count, startInner, stopInner)
		go a.RunDynamic(*monitorIntervalPtr, ch, *consumer.Ch, c, *setPointPtr, &count, startInner, stopInner)
		//go a.RunZiegler(*monitorIntervalPtr, ch, *consumer.Ch, c, *setPointPtr, &count, startInner, stopInner) // Training
	} else {
		//go a.RunTraining(*monitorIntervalPtr, ch, *consumer.Ch, c, *setPointPtr, &count, startInner, stopInner) // Training
	}

	fmt.Println("Consumer started [", *isAdaptivePtr, *controllerTypePtr, "Kp=", *kpPtr, "Ki=", *kiPtr, "Kd=", *kdPtr, "Goal=", *setPointPtr, "Monitor Interval=", *monitorIntervalPtr, "PC=", *prefetchCountPtr, "]")

	// run consumer
	consumer.Run(startInner, stopInner, startOuter, stopOuter)
}

func (c Consumer) Run(startInner chan bool, stopInner chan bool, startOuter chan bool, stopOuter chan bool) {

	for {
		select {
		case <-startOuter: // adjustment mechanism
			<-stopOuter
		case <-stopInner: // controller
			<-startInner
		default:
			d := <-c.Msgs
			d.Ack(false)                // send ack to broker
			atomic.AddUint64(&count, 1) // increment number of received messages
		}
	}
}

func (c *Consumer) configureRabbitMQ(pc int) {
	err := error(nil)

	// create connection
	c.Conn, err = amqp.Dial("amqp://guest:guest@172.26.144.1:5672/") // Leuven
	//c.Conn, err = amqp.Dial("amqp://guest:guest@192.168.0.110:5672/") // Home
	//s.Conn, err = amqp.Dial("amqp://guest:guest@172.22.38.75:5672/") // Home

	shared.FailOnError(err, "Failed to connect to RabbitMQ")

	//connSub, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	//s.ConnSub, err = amqp.Dial("amqp://nsr:nsr@localhost:5672/") // Docker 'some-rabbit'
	//shared.FailOnError(err, "Failed to connect to RabbitMQ - Subscriber")
	//defer conn.Close()

	// create channel
	c.Ch, err = c.Conn.Channel()
	shared.FailOnError(err, "Failed to open a channel")

	// declare queues
	c.Queue, err = c.Ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	shared.FailOnError(err, "Failed to declare a queue")

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
	shared.FailOnError(err, "Failed to register a consumer")

	// configure initial QoS of Req channel
	err = c.Ch.Qos(
		int(pc), // prefetch count
		0,       // prefetch size
		true,    // global TODO default is false
	)
	shared.FailOnError(err, "Failed to set QoS")
	return
}

func NewConsumer(isAdaptive bool, monitorInterval int, c controller.IController, prefetchCount int, r float64) Consumer {
	s := Consumer{}

	// Configure consumer
	s.IsAdaptive = isAdaptive
	//s.MonitorInterval = time.Duration(monitorInterval) * time.Second

	// Initialise channel to communicate with Monitor
	//s.Sensor = make(chan int, 1)

	// create Monitor
	//s.Mnt = monitor.NewMonitor(s.MonitorInterval)

	// set controller
	//s.Ctler = c

	// set goal
	//s.R = r

	// set initial PC -- always 1
	//s.PC = prefetchCount

	return s
}
*/
