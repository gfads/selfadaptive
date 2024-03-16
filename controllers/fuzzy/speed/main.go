package main

import (
	"fmt"
	"math/rand"
	"time"
)

// fuzzy sets
type FuzzySet struct {
	Name   string
	Center float64
	Width  float64
}

// fuzzy input variables
var fromEnv = FuzzySet{Name: "FromEnvironment", Center: 50, Width: 25}

// fuzzy output variable
var delay = FuzzySet{Name: "Print delay", Center: 500, Width: 400}

// fuzzy rules
func fuzzyRule(input float64) float64 {
	// Rule 1: If input is low, print delay is low
	if input <= fromEnv.Center-fromEnv.Width {
		return delay.Center - delay.Width
	}
	// Rule 2: If input is medium, print delay is medium
	if input >= fromEnv.Center-fromEnv.Width && input <= fromEnv.Center+fromEnv.Width {
		return delay.Center
	}
	// Rule 3: If input is high, print delay is high
	if input >= fromEnv.Center+fromEnv.Width {
		return delay.Center + delay.Width
	}
	return 0 // Default
}

// Fuzzy inference
func fuzzyInference(input float64) float64 {
	// Apply fuzzy rules
	output := fuzzyRule(input)
	return output
}

func main() {

	chanTimer := make(chan int)

	for i := 0; i < 20; i++ {
		// Simulate environment events
		fromEnv := float64(rand.Intn(100))

		// Fuzzy inference to determine print delay
		printDelay := fuzzyInference(fromEnv)
		//fmt.Printf("FromEnv: %.2f printDelay: %.2f\n", fromEnv, printDelay)
		fmt.Println(printDelay)

		// timer to trigger adaptation
		go myTimer(chanTimer)

		run := true
		for run == true {
			select {
			case <-chanTimer:
				run = false
			default:
				//print("+")
				time.Sleep(time.Duration(printDelay) * time.Millisecond)
			}
		}
	}
}

func myTimer(ch chan int) {
	time.Sleep(5 * time.Second)

	ch <- 0

	return
}
