package main

import (
	"fmt"
	"math/rand"
	"time"
)

const NUMBER_OF_COLORS = 7

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
		ch <- rand.Intn(NUMBER_OF_COLORS)
	}
}

func ManagingSystem(ch1, ch2 chan int) {

	for {
		t := <-ch1
		ch2 <- t
		time.Sleep(5 * time.Second)
	}
}

func ManagedSystem(ch chan int) {
	//colorReset := "\033[0m"
	b := []string{"\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m"}
	i := 0

	for {
		select {
		case i = <-ch:
		default:
		}
		fmt.Print(string(b[i]), "+")
		time.Sleep(time.Millisecond * 100)
	}
}
