package main

import (
	"fmt"
	"selfadaptive/shared"
)

func Behaviour() {
	for i := 0; i < 100; i++ {
		fmt.Print("P0")
	}
	fmt.Print(shared.ColorReset)
}
