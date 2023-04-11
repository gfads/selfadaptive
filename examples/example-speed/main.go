package main

import (
	"fmt"
	"math/rand"
	"selfadaptive/shared"
	"time"
)

func main() {

	ch1 := make(chan time.Duration)
	ch2 := make(chan time.Duration)

	go Environment(ch1)
	go ManagedSystem(ch2)
	go ManagingSystem(ch1, ch2)

	fmt.Scanln()
}

func Environment(ch chan time.Duration) {
	for {
		ch <- time.Millisecond * time.Duration(rand.Intn(1000))
	}
}

func ManagingSystem(ch1, ch2 chan time.Duration) {

	for {
		t := <-ch1
		ch2 <- t
		time.Sleep(shared.MonitorTime * time.Millisecond)
	}
}

func ManagedSystem(ch chan time.Duration) {

	t := time.Duration(0)
	for {
		select {
		case t = <-ch:
			//fmt.Println(t.Milliseconds())
		default:
		}
		fmt.Print("+")
		time.Sleep(t)
	}
}
