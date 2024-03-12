package main

import (
	"fmt"
	"main.go/shared"
	"math/rand"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go Environment(ch1)
	go ManagedSystem(ch2)
	go ManagingSystem(ch1, ch2)

	_, err := fmt.Scanln()
	if err != nil {
		return
	}
}

func Environment(ch chan int) {
	for {
		ch <- rand.Intn(shared.NumberOfColors)
	}
}

func ManagingSystem(ch1, ch2 chan int) {

	for {
		t := <-ch1
		ch2 <- t
		time.Sleep(shared.MonitorTime * time.Second)
	}
}

func ManagedSystem(ch chan int) {

	i := 0
	for {
		select {
		case i = <-ch:
		default:
		}
		fmt.Print(shared.ColorBehaviours[i], "+")
		time.Sleep(time.Millisecond * 100)
	}
}
