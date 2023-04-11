package main

import (
	"fmt"
	"math/rand"
	"selfadaptive/controllers/onoff"
	"time"
)

func main() {

	ch := make(chan time.Duration)

	go ManagedSystem(ch)
	go ManagingSystem(ch)

	fmt.Scanln()
}

func ManagingSystem(ch chan time.Duration) {
	c := onoff.OnOff{}

	for {
		t := <-ch
		ch <- c.Update(t, 590*time.Millisecond)
		time.Sleep(5000 * time.Millisecond)
	}
}

func ManagedSystem(ch chan time.Duration) {

	t := time.Millisecond * time.Duration(rand.Intn(1000))

	for {
		select {
		case ch <- t:
			t = <-ch
			//fmt.Println(t.Milliseconds())
		default:
		}
		fmt.Print("+")
		time.Sleep(t)
	}
}
