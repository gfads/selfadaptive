package main

import (
	"fmt"
	"math"
	"math/rand"
)

const NB = -3.0
const NM = -2.0
const PB = 3.0
const PM = 2.0
const Zero = 0.0

func main() {
	var errorSet = []float64{-50, -25, 0, 25, 50}
	var pcMoveSet = []float64{-5, 0, 5} // increase, nothing, decrease

	goal := 400.0

	pc := 0.0

	for i := 0; i < 10; i++ {
		rate := float64(myRandon(250, 450))

		// Fuzzification
		e := goal - rate // error
		inputFuzzy := fuzzification(errorSet, e)

		// Fuzzy rules
		outputFuzzy := inputFuzzy

		// Defuzzification
		pc += centroidDefuzzification(pcMoveSet, outputFuzzy)

		fmt.Printf("[%2.f msg/s, %2.f]\n", rate, pc)
	}
}

// Define membership functions
func triangularMF(x float64, a float64, b float64, c float64) float64 {
	return math.Max(0, math.Min((math.Min((x-a)/(b-a), (c-x)/(c-b))), 1))
}

// Fuzzification: Fuzzy input mapping
func fuzzification(errorSet []float64, error float64) []float64 {
	r := []float64{}

	for i := 0; i < len(errorSet)-2; i++ {
		r = append(r, triangularMF(error, errorSet[i], errorSet[i+1], errorSet[i+2]))
	}
	return r
}

// Defuzzification: Centroid method
func centroidDefuzzification(pcMoveSet []float64, outputFuzzy []float64) float64 {
	r := 0.0

	numerator := 0.0
	for i := 0; i < len(pcMoveSet); i++ {
		numerator += outputFuzzy[i] * pcMoveSet[i]
	}
	denominator := 0.0
	for i := 0; i < len(pcMoveSet); i++ {
		denominator += outputFuzzy[i]
	}

	if denominator == 0 {
		r = 1.0 // minimum PC
	} else {
		r = numerator / denominator
	}

	if math.Round(r) == 0 {
		r = 1
	}

	return math.Round(r)
}

func myRandon(min, max int) float64 {
	for {
		n := rand.Intn(max)
		if n >= min {
			return float64(n)
		}
	}
}
