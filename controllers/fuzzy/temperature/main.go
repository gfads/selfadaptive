package main

import (
	"fmt"
	"math/rand"
)

// Define fuzzy sets
type FuzzySet struct {
	Name   string
	Center float64
	Width  float64
}

// Define fuzzy input variables
var temperature = FuzzySet{Name: "Temperature", Center: 25, Width: 5}

// Define fuzzy output variable
var fanSpeed = FuzzySet{Name: "FanSpeed", Center: 0, Width: 10}

// Define fuzzy rules
func fuzzyRule(input float64) float64 {
	// Rule 1: If temperature is low, fan speed is low
	if input <= temperature.Center-temperature.Width {
		return 0
	}
	// Rule 2: If temperature is medium, fan speed is medium
	if input >= temperature.Center-temperature.Width && input <= temperature.Center+temperature.Width {
		return 5
	}
	// Rule 3: If temperature is high, fan speed is high
	if input >= temperature.Center+temperature.Width {
		return 10
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
	for i := 1; i < 10; i++ {
		// Simulated temperature input
		currentTemperature := float64(rand.Intn(60))

		// Fuzzy inference to determine fan speed
		fanSpeedOutput := fuzzyInference(currentTemperature)

		fmt.Printf("Temperature: %.2f Fan speed: %.2f\n", currentTemperature, fanSpeedOutput)
	}
}
