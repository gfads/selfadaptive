package main

import (
	"fmt"
	"main.go/controllers/fuzzy"
	"time"
)

var Setpoint int = 500

func main() {

	ch := make(chan float64)

	go ManagedSystem(ch)
	go ManagingSystem(ch)

	fmt.Scanln()
}

func ManagingSystem(ch chan float64) {
	//c := onoffbasic.Controller{}
	//c.Initialise(-1, 400, 600)
	//c := basicpid.Controller{}
	//c.Initialise(-1, 100, 1000, 2, 1, 0)
	//c := hpa.Controller{}
	//c.Initialise(-1, 400, 600, 100)
	c := fuzzy.Controller{}
	c.Initialise(500, 100, 500, 100)

	chTimer := make(chan int)
	for i := 0; i < 20; i++ {
		go myTimer(chTimer)

		select {
		case <-chTimer:
			ch <- 0
			y := <-ch
			//ch <- c.Update(float64(Setpoint), float64(y)) // no fuzzy
			ch <- c.Update(float64(y)) // fuzzy
		}
		//time.Sleep(5000 * time.Millisecond)
	}
}

func ManagedSystem(ch chan float64) {
	var u float64

	env := []float64{81, 887, 847, 59, 81, 318, 425, 540, 456, 300, 694, 511, 162, 89, 728, 274, 211, 445, 237, 106}
	i := 0

	for {
		select {
		case <-ch:
			//y := rand.Intn(1000)
			//println(y)
			y := env[i]
			ch <- y
			u = <-ch
			fmt.Println(u)
			i++
		default:
		}
		//fmt.Print("+")
		time.Sleep(time.Millisecond * time.Duration(u))
	}
}

func myTimer(ch chan int) {
	time.Sleep(1000 * time.Millisecond)

	ch <- 0

	return
}
