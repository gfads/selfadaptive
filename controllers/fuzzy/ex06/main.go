package main

import (
	"fmt"
)

// Linguistic terms for error and change in error
const (
	NegativeBig   = "NegativeBig"
	NegativeSmall = "NegativeSmall"
	Zero          = "Zero"
	PositiveSmall = "PositiveSmall"
	PositiveBig   = "PositiveBig"

	DecreasingFast = "DecreasingFast"
	DecreasingSlow = "DecreasingSlow"
	Stable         = "Stable"
	IncreasingSlow = "IncreasingSlow"
	IncreasingFast = "IncreasingFast"
)

// Define membership functions for error
func membershipError(errorValue float64) map[string]float64 {
	return map[string]float64{
		NegativeBig:   max((0-errorValue)/20, 0),
		NegativeSmall: max(errorValue/20, min((0-errorValue)/20, (40-errorValue)/20)),
		Zero:          max(min((errorValue-20)/20, (60-errorValue)/20), 0),
		PositiveSmall: max((errorValue-40)/20, min((errorValue-20)/20, (100-errorValue)/20)),
		PositiveBig:   max((errorValue-60)/20, 0),
	}
}

// Define membership functions for change in error
func membershipChangeInError(changeInError float64) map[string]float64 {
	return map[string]float64{
		DecreasingFast: max(min((0-changeInError)/10, 1), 0),
		DecreasingSlow: max(min(changeInError/10, (20-changeInError)/20), 0),
		Stable:         max(min((changeInError-10)/10, (30-changeInError)/20), 0),
		IncreasingSlow: max(min(changeInError/10, (40-changeInError)/20), 0),
		IncreasingFast: max((changeInError-20)/10, 0),
	}
}

// Rule evaluation
func ruleEvaluation(errorValues, changeValues map[string]float64) map[string]float64 {
	powerLevel := make(map[string]float64)

	// Define fuzzy rules for the power level based on error and change in error
	powerLevel["Low"] = min(errorValues[NegativeBig], changeValues[IncreasingFast])
	powerLevel["Medium"] = min(errorValues[Zero], changeValues[Stable])
	powerLevel["High"] = min(errorValues[PositiveBig], changeValues[DecreasingFast])

	return powerLevel
}

// Defuzzify power level using centroid method
func defuzzify(powerLevel map[string]float64) float64 {
	numerator, denominator := 0.0, 0.0
	for key, value := range powerLevel {
		switch key {
		case "Low":
			numerator += value * 10 // arbitrary value for low power
		case "Medium":
			numerator += value * 50 // arbitrary value for medium power
		case "High":
			numerator += value * 90 // arbitrary value for high power
		}
		denominator += value
	}
	return numerator / denominator
}

// Helper functions
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	// Example crisp error and change in error values
	crispError := 30.0
	crispChange := 5.0

	// Fuzzification for error and change in error
	errorValues := membershipError(crispError)
	changeValues := membershipChangeInError(crispChange)

	// Rule evaluation
	powerLevel := ruleEvaluation(errorValues, changeValues)

	// Defuzzification
	crispOutput := defuzzify(powerLevel)

	fmt.Printf("Crisp Power Level: %.2f\n", crispOutput)
}
