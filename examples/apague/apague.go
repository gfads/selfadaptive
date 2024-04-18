package main

import "fmt"

func main() {

	ch := make(chan bool, 2)

	ch <- true
	ch <- false

	for m := range ch {
		fmt.Println(m)
	}
}
