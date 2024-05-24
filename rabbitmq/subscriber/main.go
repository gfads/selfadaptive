package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"main.go/controllers/def/parameters"
	"main.go/rabbitmq/adaptationlogic"
	"main.go/rabbitmq/mytimer"
	"main.go/shared"
	_ "net/http/pprof"
	"os"
	"time"
)

type Subscriber struct {
	ExecutionType string
	Conn          *amqp.Connection
	Ch            *amqp.Channel
	Queue         amqp.Queue
	Msgs          <-chan amqp.Delivery
	PC            int
}

func main() {

	// load parameters
	e := parameters.ExecutionParameters{}
	p := e.Load()

	// validate parameters
	e.Validate(p) // TODO

	// show parameters
	e.Show(p)

	// create new consumer
	var consumer = NewConsumer(*p.ExecutionType, *p.PrefetchCount)

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
	toAdapter := make(chan shared.SubscriberToAdapter) // no. of messages
	fromAdapter := make(chan int)                      // pc
	startTimer := make(chan bool)                      // start timer
	stopTimer := make(chan bool)                       // stop timer

	// define and open csv file to record experiment results
	dataFileName := *p.OutputFile
	df, err := os.Create(shared.DockerDir + "/" + dataFileName + ".csv")
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// create & start timer
	t := mytimer.NewMyTimer(*p.MonitorInterval, startTimer, stopTimer)
	go t.RunMyTimer()

	// create adapter
	adapter := adaptationlogic.NewAdaptationLogic(toAdapter, fromAdapter, p, df)

	consumer.Warmup()

	if *p.ExecutionType == shared.OpenLoop {
		consumer.RunOpenLoop(startTimer, stopTimer, p, df)
	} else {
		// stard adapter
		go adapter.Run()

		// start consumer
		consumer.RunClosedLoop(startTimer, stopTimer, toAdapter, fromAdapter, df)
	}
}

func (c Subscriber) RunOpenLoop(startTimer, stopTimer chan bool, p parameters.ExecutionParameters, df *os.File) {
	count := 0
	nSameLevel := 0
	c.PC = *p.PrefetchCount

	var t1, t2 time.Time

	count = 0 // initialise no. of received messages
	for d := range c.Msgs {
		err := d.Ack(false) // send ack to broker
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		count++ // increment number of received messages
		select {
		case <-startTimer: // start monitor timer
			t1 = time.Now()
		case <-stopTimer: // stop monitor timer
			t2 = time.Now()
			nSameLevel++
			rate1 := float64(count) / float64(*p.MonitorInterval)
			rate2 := float64(count) / float64(t2.Sub(t1).Seconds())
			count = 0 // re-initialise no. of received messages

			fmt.Printf("%v %v \n", rate1, rate2)
			//fmt.Fprintf(df, "%d;%d;%f;%f\n", 0, c.PC, rate, 0.0) // queue size; pc;rate;goal
			//fmt.Printf("%d;%d;%f;%f\n", 0, c.PC, rate, 0.0)
			if nSameLevel > shared.SizeOfSameLevel {
				c.PC++
				nSameLevel = 0
				if c.PC >= shared.TrainingSampleSize {
					fmt.Println("End of Experiment!!")
					df.Close()
					os.Exit(0)
				} else {
					// configure new pc
					err := c.Ch.Qos(
						c.PC, // prefetch count
						0,    // prefetch size
						true, // global TODO default is false
					)
					if err != nil {
						shared.ErrorHandler(shared.GetFunction(), "Failed to set QoS")
					}
				}
			}
		default:
		}
	}
}

func (c Subscriber) Warmup() {

	fmt.Println("Begin of Warming up...")

	// configure pc to zero
	err := c.Ch.Qos(
		0,    // prefetch count
		0,    // prefetch size
		true, // global TODO default is false
	)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to set QoS")
	}
	for i := 0; i < 1000000; i++ {
		//for i := 0; i < 10; i++ {
		d := <-c.Msgs
		err := d.Ack(false) // send ack to broker
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}

	// return pc back to initial pc
	err = c.Ch.Qos(
		c.PC, // prefetch count
		0,    // prefetch size
		true, // global TODO default is false
	)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to set QoS")
	}
	fmt.Println("End of Warming up...")
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

func (c Subscriber) RunNonAdaptiveMonitored(startTimer, stopTimer chan bool, p parameters.ExecutionParameters) {
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

func (c Subscriber) RunClosedLoop(startTimer, stopTimer chan bool, toAdapter chan shared.SubscriberToAdapter, fromAdapter chan int, f *os.File) {

	count := 0
	var t1, t2 time.Time

	for d := range c.Msgs {
		err := d.Ack(false) // send ack to broker
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		count++ // increment number of received messages
		select {
		case <-startTimer: // start monitor timer
			t1 = time.Now()
		case <-stopTimer: // stop monitor timer
			t2 = time.Now()
			// inspect queue
			q, err1 := c.Ch.QueueInspect("rpc_queue")
			if err1 != nil {
				shared.ErrorHandler(shared.GetFunction(), "Impossible to inspect the queue")
			}

			// send no. of received messages to adaptation logic
			toAdapter <- shared.SubscriberToAdapter{ReceivedMessages: count, QueueSize: q.Messages, D: t2.Sub(t1).Seconds()}

			// re-initialise message counter
			count = 0

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
		default:
		}
	}
}

func (c *Subscriber) ConfigureRabbitMQ(pc int) {
	err := error(nil)

	// create connection
	c.Conn, err = amqp.Dial("amqp://guest:guest@" + shared.IpPortRabbitMQ + "/")

	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to connect to RabbitMQ")
	}

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

func NewConsumer(e string, pc int) Subscriber {
	s := Subscriber{ExecutionType: e, PC: pc}

	return s
}
