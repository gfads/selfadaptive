package main

import (
	"fmt"
	"math"
)

const LP = "LP" // Large Positive
const SP = "SP" // Small Positive
const ZE = "ZE" // Zero
const SN = "SN" // Smal Negative
const LN = "LN" // Large Negative

func main() {
	//fuzzySet := []string{LP, SP, ZE, SN, LN}

	errorExamples := []float64{0.9}
	deltaExamples := []float64{0.50}

	for i := 0; i < len(errorExamples); i++ {
		e := errorExamples[i]
		d := deltaExamples[i]

		// fuzzyfication
		fuzzifiedSetError := fuzzyfication(e)
		fuzzifiedSetDelta := fuzzyfication(d)

		// apply rules
		mx, output := applyRules(fuzzifiedSetError, fuzzifiedSetDelta)

		// Deffuzification
		u := centroidDeffuzification(mx, output)
		fmt.Printf("Error: %.2f Delta=%.2f u = %.2f \n", e, d, u)
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

func applyRules(e map[string]float64, d map[string]float64) ([]float64, []float64) {

	mx := []float64{}
	output := []float64{}

	// Rule 1:  IF e = ZE AND delta = ZE THEN output = ZE
	eR := e[ZE]
	dR := d[ZE]
	m1 := math.Min(eR, dR)
	o1 := getMax(ZE)
	mx = append(mx, m1)
	output = append(output, o1)

	// Rule 2:  IF e = ZE AND delta = SP THEN output = SN
	eR = e[ZE]
	dR = d[SP]
	m2 := math.Min(eR, dR)
	o2 := getMax(SN)
	mx = append(mx, m2)
	output = append(output, o2)

	// Rule 3:  IF e = SN AND delta = SN THEN output = LP
	eR = e[SN]
	dR = d[SN]
	m3 := math.Min(eR, dR)
	o3 := getMax(LP)
	mx = append(mx, m3)
	output = append(output, o3)

	// Rule 4:  IF e = LP OR  delta = LP THEN output = LN
	eR = e[LP]
	dR = d[LP]
	m4 := math.Max(eR, dR)
	o4 := getMax(LN)
	mx = append(mx, m4)
	output = append(output, o4)

	return mx, output
}

func fuzzyfication(n float64) map[string]float64 {

	r := map[string]float64{}

	r[LP] = triangularMF(n, 0.5, 0.75, 1.0)
	r[SP] = triangularMF(n, 0.0, 0.5, 1.0)
	r[ZE] = triangularMF(n, -0.5, 0.0, 0.5)
	r[SN] = triangularMF(n, -1.0, -0.5, 0.0)
	r[LN] = triangularMF(n, -1.0, -0.75, -0.5)
	return r
}

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
