package main

import (
	"fmt"
	"math"
	"selfadaptive/shared"
)

func Behaviour() {
	for i := 0; i < 100; i++ {
		fmt.Print("Agora estou calculando o seno de um número: ", math.Sin(10.0))
	}
	fmt.Print(shared.ColorReset)
}
