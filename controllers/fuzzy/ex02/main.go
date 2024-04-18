package main

import (
	"fmt"
	"math"
)

// Define membership functions
func triangularMF(x float64, a float64, b float64, c float64) float64 {
	return math.Max(0, math.Min((math.Min((x-a)/(b-a), (c-x)/(c-b))), 1))
}

// Fuzzification: Fuzzy input mapping
func fuzzification(temperature float64) (float64, float64, float64) {
	cold := triangularMF(temperature, 0, 10, 20)
	moderate := triangularMF(temperature, 10, 20, 30)
	hot := triangularMF(temperature, 20, 30, 40)

	return cold, moderate, hot
}

// Defuzzification: Centroid method
func defuzzification(heaterIncrease float64, heaterMaintain float64, heaterDecrease float64) float64 {
	numerator := (heaterIncrease * 10) + (heaterMaintain * 5) + (heaterDecrease * 0)
	denominator := heaterIncrease + heaterMaintain + heaterDecrease

	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}

func main() {
	// Input set (Temperature) = {0, 10, 20, 30, 40}
	// Cold = TriangularMF (0, 10, 20)
	// Moderate = TriangularMF (10, 20, 30)
	// Hot = TriangularMF(20, 30, 40)

	// Output set (Heater) = {0, 5, 10}
	// heaterIncrease = Cold
	// heaterMaintain = Moderate
	// heaterDecrease = Hot
	// Centroid

	// Test the fuzzy controller with a temperature value
	temperature := 15.0

	// Fuzzification
	cold, moderate, hot := fuzzification(temperature)

	// Fuzzy Rules
	heaterIncrease := cold
	heaterMaintain := moderate
	heaterDecrease := hot

	// Defuzzification
	heaterOutput := defuzzification(heaterIncrease, heaterMaintain, heaterDecrease)

	fmt.Printf("Temperature: %.2f°C\n", temperature)
	fmt.Printf("Cold: %.2f, Moderate: %.2f, Hot: %.2f\n", cold, moderate, hot)
	fmt.Printf("Heater Output: %.2f\n", heaterOutput)
}
