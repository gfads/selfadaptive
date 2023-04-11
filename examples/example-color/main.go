package main

import (
	"fmt"
	"math/rand"
	"selfadaptive/shared"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go Environment(ch1)
	go ManagedSystem(ch2)
	go ManagingSystem(ch1, ch2)

	fmt.Scanln()
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
		fmt.Print(string(shared.ColorBehaviours[i]), "+")
		time.Sleep(time.Millisecond * 100)
	}
}
