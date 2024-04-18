package main

import (
	"fmt"
	"math"
)

// Define triangular membership function
func triangularMF(x float64, a float64, b float64, c float64) float64 {
	return math.Max(0, math.Min((math.Min((x-a)/(b-a), (c-x)/(c-b))), 1))
}

// Define fuzzy rules function
func fuzzyRules(temperature float64) float64 {
	// Define membership functions for temperature
	cold := triangularMF(temperature, 0, 10, 20)
	moderate := triangularMF(temperature, 10, 20, 30)
	hot := triangularMF(temperature, 20, 30, 40)

	// Define fuzzy rules
	// Rule 1: If the temperature is cold, then increase the heater
	// Rule 2: If the temperature is moderate, then maintain the heater
	// Rule 3: If the temperature is hot, then decrease the heater
	heaterIncrease := cold
	heaterMaintain := moderate
	heaterDecrease := hot

	// Aggregate the fuzzy outputs
	aggregatedOutput := (heaterIncrease*10 + heaterMaintain*5 + heaterDecrease*0) /
		(heaterIncrease + heaterMaintain + heaterDecrease)

	return aggregatedOutput
}

func main() {
	// Test the fuzzy rules with a temperature value
	temperature := 25.0
	heaterOutput := fuzzyRules(temperature)
	fmt.Printf("Temperature: %.2f°C, Heater Output: %.2f\n", temperature, heaterOutput)
}
