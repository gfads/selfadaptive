package main

import (
	"fmt"
	"math"
)

// Define membership functions
func triangularMF(x float64, a float64, b float64, c float64) float64 {
	return math.Max(0, math.Min((x-a)/(b-a), (c-x)/(c-b)))
}

func gaussianMF(x float64, m float64, d float64) float64 {
	temp := -math.Pow(x-m, 2.0) / 2.0 * math.Pow(d, 2.0)
	r := math.Exp(temp)

	return r
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

func trapezoidalMF(x float64, a float64, b float64, c float64, d float64) float64 {
	min1 := math.Min(x-a/b-a, 1.0)
	r := math.Max(math.Min(min1, d-x/d-c), 0.0)

	return r
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

	for i := 0; i > -120; i -= 1 {
		fmt.Printf("%v %.2f\n", i, gaussianMF(float64(i), -100.0, -0.1))
	}
}
