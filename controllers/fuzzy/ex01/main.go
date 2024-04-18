package main

import (
	"fmt"
	"math"
)

// Define membership functions
func triangularMF(x float64, a float64, b float64, c float64) float64 {
	return math.Max(0, math.Min((x-a)/(b-a), (c-x)/(c-b)))
}

// Rectangular membership function
func rectangularMF(x float64, a float64, b float64) float64 {
	if x >= a && x <= b {
		return 1
	} else {
		return 0
	}
}

func rampMF(x float64, a float64, b float64) float64 {
	if x <= a {
		return 0
	} else if x >= b {
		return 1
	} else {
		return (x - a) / (b - a)
	}
}

// Pi-shaped membership function
func piMF(x float64, a float64, b float64, c float64, d float64) float64 {
	if x < a || x > d {
		return 0
	} else if x >= a && x <= b {
		return (x - a) / (b - a)
	} else if x >= c && x <= d {
		return (d - x) / (d - c)
	} else {
		return 1
	}
}

// Trapezoidal membership function
func trapezoidalMF(x float64, a float64, b float64, c float64, d float64) float64 {
	if x < a || x > d {
		return 0
	} else if x >= b && x <= c {
		return 1
	} else if x >= a && x < b {
		return (x - a) / (b - a)
	} else if x > c && x <= d {
		return (d - x) / (d - c)
	} else {
		return 1
	}
}

// Define fuzzy rules
func fuzzyRules(input float64) float64 {
	if input < 5 {
		return 0
	} else if input >= 5 && input < 10 {
		//return triangularMF(input, 5, 7.5, 10)
		//return rampMF(input, 5, 10)
		return rectangularMF(input, 5, 10)
	} else if input >= 10 && input < 15 {
		//return triangularMF(input, 10, 12.5, 15)
		//return rampMF(input, 10, 15)
		return rectangularMF(input, 10, 15)
	} else {
		return 1
	}
}

func main() {
	// Test the fuzzy controller
	input := 7.8
	output := fuzzyRules(input)
	fmt.Printf("Input: %f, Output: %f\n", input, output)
}
