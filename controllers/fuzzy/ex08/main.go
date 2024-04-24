package main

import (
	"fmt"
	"math"
)

const NB = "NB" // Negative Big
const NM = "NM" // Negative Medium
const PB = "PB" // Positive Big
const PM = "PM" // Positive Medium

const POS = "Pos"
const NEG = "Neg"
const ZERO = "Zero"

func main() {
	//fuzzySet := []string{LP, SP, ZE, SN, LN}

	// Example 03 - 2005-Design Of Fuzzy Controllers.pdf

	errorExamples := []float64{-50}
	changeExamples := []float64{0.50}

	for i := 0; i < len(errorExamples); i++ {
		e := errorExamples[i]
		c := changeExamples[i]

		// fuzzyfication
		fuzzifiedSetError := fuzzyfication(e)
		fuzzifiedSetChange := fuzzyfication(c)

		// apply rules
		mx, output := applyRules(e, c, fuzzifiedSetError, fuzzifiedSetChange)

		// Deffuzification
		u := centroidDeffuzification(mx, output)
		fmt.Printf("Error: %.2f u = %.2f \n", e, u)
	}
}

func centroidDeffuzification(mx, output []float64) float64 {

	numerator := 0.0
	denominator := 0.0

	for i := 0; i < len(mx); i++ {
		numerator = numerator + mx[i]*output[i]
		denominator = denominator + mx[i]
	}
	u := 0.0
	if denominator == 0 {
		u = 1
	} else {
		u = numerator / denominator
	}
	return u
}

func applyRules(e float64, c float64, ev map[string][]float64, cv map[string][]float64) ([]float64, []float64) {

	mxTemp := [][]float64{}
	output := []float64{}

	//1. Rule 1 - Figure 8
	v := ev[NEG]
	ef := cosineMF(-100, -100, -60, 10, e)
	temp := []float64{}
	for i := 0; i < len(v); i++ {
		temp = append(temp, math.Min(ef, v[i]))
	}
	mxTemp = append(mxTemp, temp)

	//2. Rule 2 - Figure 8
	v = ev[ZERO]
	ef = cosineMF(-90, -20, 20, 90, e)
	temp = []float64{}
	for i := 0; i < len(v); i++ {
		temp = append(temp, math.Min(ef, v[i]))
	}
	mxTemp = append(mxTemp, temp)

	//3. Rule 3 - Figure 8
	v = ev[POS]
	ef = cosineMF(-10, -60, 100, 100, e)
	temp = []float64{}
	for i := 0; i < len(v); i++ {
		temp = append(temp, math.Min(ef, v[i]))
	}
	mxTemp = append(mxTemp, temp)

	v = []float64{}
	for i := 0; i < len(mxTemp[0]); i++ {
		max := -10000.0
		for j := 0; j < len(mxTemp); j++ {
			if mxTemp[j][i] > max {
				max = mxTemp[j][i]
			}
		}
		v = append(v, max)
	}

	mx := v
	output = []float64{-100, -50, 0, 50, 100}

	return mx, output
}

func fuzzyfication(n float64) map[string][]float64 {

	r := map[string][]float64{}

	// NEG
	v := []float64{}
	v = append(v, cosineMF(-100, -100, -60, 10, -100))
	v = append(v, cosineMF(-100, -100, -60, 10, -50))
	v = append(v, cosineMF(-100, -100, -60, 10, 0))
	v = append(v, cosineMF(-100, -100, -60, 10, 50))
	v = append(v, cosineMF(-100, -100, -60, 10, 100))
	r[NEG] = v

	// ZERO
	v = []float64{}
	v = append(v, cosineMF(-90, -20, 20, 90, -100))
	v = append(v, cosineMF(-90, -20, 20, 90, -50))
	v = append(v, cosineMF(-90, -20, 20, 90, 0))
	v = append(v, cosineMF(-90, -20, 20, 90, 50))
	v = append(v, cosineMF(-90, -20, 20, 90, 100))
	r[ZERO] = v

	// POS
	v = []float64{}
	v = append(v, cosineMF(-10, 60, 100, 100, -100))
	v = append(v, cosineMF(-10, 60, 100, 100, -50))
	v = append(v, cosineMF(-10, 60, 100, 100, 0))
	v = append(v, cosineMF(-10, 60, 100, 100, 50))
	v = append(v, cosineMF(-10, 60, 100, 100, 100))
	r[POS] = v

	return r
}

/*
func getMax(s string) float64 {
	r := 0.0
	max := -1000.0

	for i := -1.0; i <= 1.0; i += 0.25 {
		v := fuzzyfication(i)
		if v[s] > max {
			max = v[s]
			r = i
		}
	}
	return r
}
*/
func triangularMF(x float64, a float64, b float64, c float64) float64 {
	return math.Max(0, math.Min((x-a)/(b-a), (c-x)/(c-b)))
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

func cosineMF(x1, x2, x3, x4, x float64) float64 {
	r := 0.0

	// s-curve
	rs := 0.0
	if x2-x1 == 0 { // TODO
		rs = 10000.00
	} else {
		if x < x1 {
			rs = 0.0
		} else if x1 <= x && x <= x2 {
			rs = 1.0/2.0 + 1.0/2.0*math.Cos(((x-x2)/(x2-x1))*math.Pi)
		} else {
			rs = 1.0
		}
	}

	// z-curve
	rz := 0.0
	if (x4 - x3) == 0.0 {
		rz = 10000.00
	} else {
		if x < x3 {
			rz = 1.0
		} else if x3 <= x && x <= x4 {
			rz = 1.0/2.0 + 1.0/2.0*math.Cos((x-x3)/(x4-x3)*math.Pi)
		} else {
			rz = 0.0
		}
	}
	r = math.Min(rs, rz)

	return r
}

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
