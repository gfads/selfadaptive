package main

import (
	"fmt"
	"math"
)

const LARGEPOSITIVE = "LP" // Large Positive
const SMALLPOSITIVE = "SP" // Small Positive
const ZERO = "ZE"          // Zero
const SMALLNEGATIVE = "SN" // Smal Negative
const LARGENEGATIVE = "LN" // Large Negative

const LARGEINCREASE = "LI"  // Large Positive
const SMALLINCREASE = "SI"  // Small Positive
const MAINTAIN = "MAINTAIN" // Zero
const SMALLDECREASE = "SD"  // Small Negative
const LARGEDECREASE = "LD"  // Large Negative

func main() {
	//fuzzySet := []string{LP, SP, ZE, SN, LN} // Input
	//fuzzySet := []string{LPPC,SPPC,ONEPC,SNPC,LNPC} // Output

	errorExamples := []float64{-90, -50, 0, 50, 90}
	//deltaExamples := []float64{50}

	for i := 0; i < len(errorExamples); i++ {
		e := errorExamples[i]
		//d := deltaExamples[i]

		// fuzzyfication
		fuzzifiedSetError := fuzzyficationInput(e)
		//fuzzifiedSetDelta := fuzzyficationInput(d)

		// apply rules
		mx, output := applyRules(fuzzifiedSetError)

		// Deffuzification
		u := centroidDeffuzification(mx, output)
		fmt.Printf("Error: %.2f ACTION ON PC = %.2f \n", e, u)
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

func applyRules(e map[string]float64) ([]float64, []float64) {

	mx := []float64{}
	output := []float64{}

	// Rule 1:  IF e = LARGEPOSITIVE THEN output = LARGEINCREASE
	eR := e[LARGEPOSITIVE]
	//dR := d[ZE]
	//m1 := math.Min(eR, dR)
	m1 := eR
	o1 := getMaxOutput(LARGEINCREASE)
	mx = append(mx, m1)
	output = append(output, o1)

	// Rule 2:  IF e = SMALLPOSITIVE THEN output = SMALLINCREASE
	eR = e[SMALLPOSITIVE]
	//dR = d[SP]
	//m2 := math.Min(eR, dR)
	m2 := eR
	o2 := getMaxOutput(SMALLINCREASE)
	mx = append(mx, m2)
	output = append(output, o2)

	// Rule 3:  IF e = ZE THEN output = MAINTAIN
	eR = e[ZERO]
	//dR = d[SN]
	//m3 := math.Min(eR, dR)
	m3 := eR
	o3 := getMaxOutput(MAINTAIN)
	mx = append(mx, m3)
	output = append(output, o3)

	// Rule 4:  IF e = SN THEN output = SMALLPC
	eR = e[SMALLNEGATIVE]
	//dR = d[LP]
	//m4 := math.Max(eR, dR)
	m4 := eR
	o4 := getMaxOutput(SMALLDECREASE)
	mx = append(mx, m4)
	output = append(output, o4)

	// Rule 5:  IF e = LN THEN output = LARGEPC
	eR = e[LARGENEGATIVE]
	//dR = d[LP]
	//m4 := math.Max(eR, dR)
	m5 := eR
	o5 := getMaxOutput(LARGEDECREASE)
	mx = append(mx, m5)
	output = append(output, o5)

	return mx, output
}

func fuzzyficationInput(n float64) map[string]float64 {

	r := map[string]float64{}

	r[LARGEPOSITIVE] = triangularMF(n, 50, 75, 100)
	r[SMALLPOSITIVE] = triangularMF(n, 0, 50, 100)
	r[ZERO] = triangularMF(n, -50, 0, 50)
	r[SMALLNEGATIVE] = triangularMF(n, -100, -50, 0)
	r[LARGENEGATIVE] = triangularMF(n, -100, -75, -50)

	//fmt.Println(r[LP], r[SP], r[ZE], r[SN], r[LN])
	return r
}

func fuzzyficationOutput(n float64) map[string]float64 {

	r := map[string]float64{}

	r[LARGEINCREASE] = triangularMF(n, 2.0, 3.0, 4.0)
	r[SMALLINCREASE] = triangularMF(n, 1.0, 2.0, 3.0)
	r[MAINTAIN] = triangularMF(n, 1.0, 0.0, -1.0)
	r[SMALLDECREASE] = triangularMF(n, -1.0, -2.0, -3.0)
	r[LARGEDECREASE] = triangularMF(n, -2.0, -3.0, -4.0)

	/*r[VERYLARGEPC] = rampMF(n, 5.0, 6.0)
	r[LARGEPC] = rampMF(n, 4.0, 5.0)
	r[MEDIUMPC] = rampMF(n, 3.0, 4.0)
	r[SMALLPC] = rampMF(n, 2.0, 3.0)
	r[ONEPC] = rampMF(n, 1.0, 2.0)
	*/
	return r
}

func getMaxOutput(s string) float64 {
	r := 0.0
	max := -10000.0

	for i := -4.0; i <= 5.0; i += 1.0 {
		v := fuzzyficationOutput(i)
		//fmt.Println(i, "->", v[ONEPC], v[SMALLPC], v[MEDIUMPC], v[LARGEPC], v[VERYLARGEPC])
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
