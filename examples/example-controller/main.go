package main

import (
	"fmt"
	"main.go/controllers/hpa"
	"math/rand"
	"time"
)

var Setpoint int = 500

func main() {

	ch := make(chan int)

	go ManagedSystem(ch)
	go ManagingSystem(ch)

	fmt.Scanln()
}

func ManagingSystem(ch chan int) {
	//c := onoffbasic.Controller{}
	//c.Initialise(-1, 100, 1000)
	//c := basicpid.Controller{}
	//c.Initialise(-1, 100, 1000, 2, 1, 0)
	c := hpa.Controller{}
	c.Initialise(-1, 100, 1000, 100)

	for {
		y := <-ch
		ch <- int(c.Update(float64(Setpoint), float64(y)))
		time.Sleep(5000 * time.Millisecond)
	}
}

func ManagedSystem(ch chan int) {
	var u int

	for {
		y := rand.Intn(1000)
		select {
		case ch <- y:
			u = <-ch
			//fmt.Println(u)
		default:
		}
		fmt.Print("+")
		time.Sleep(time.Millisecond * time.Duration(u))
	}
}
