package main

import (
	"fmt"
	"math"
	"os"
)

const LP = "LP"
const SP = "SP"
const ZE = "ZE"
const SN = "SN"
const LN = "LN"

type X struct {
	Name   string
	Input  float64
	Output float64
}

func main() {
	//FuzzySet := []string{LP, SP, ZE, SN, LN}

	// rules
	error := 0.25
	delta := 0.50

	// fuzzyfication
	table2 := []X{}
	table2 = append(table2, X{LP, fuzzyfication(LP, error), fuzzyfication(LP, delta)})
	table2 = append(table2, X{SP, fuzzyfication(SP, error), fuzzyfication(SP, delta)})
	table2 = append(table2, X{ZE, fuzzyfication(ZE, error), fuzzyfication(ZE, delta)})
	table2 = append(table2, X{SN, fuzzyfication(SN, error), fuzzyfication(SN, delta)})
	table2 = append(table2, X{LN, fuzzyfication(LN, error), fuzzyfication(LN, delta)})

	// apply rules
	fmt.Println(table2)

	mx := []float64{}
	outputs := []float64{}

	// Rule 1:  IF e = ZE AND delta = ZE THEN output = ZE
	eR := getOutputError(table2, ZE)
	dR := getOutputDelta(table2, ZE)
	m1 := math.Min(eR, dR)
	o1 := 0.0
	mx = append(mx, m1)
	outputs = append(outputs, o1)

	// Rule 2:  IF e = ZE AND delta = SP THEN output = SN
	eR = getOutputError(table2, ZE)
	dR = getOutputDelta(table2, SP)
	m2 := math.Min(eR, dR)
	o2 := -0.5
	mx = append(mx, m2)
	outputs = append(outputs, o2)

	// Rule 3:  IF e = SN AND delta = SN THEN output = LP
	eR = getOutputError(table2, SN)
	dR = getOutputDelta(table2, SN)
	m3 := math.Min(eR, dR)
	o3 := 1.0
	mx = append(mx, m3)
	outputs = append(outputs, o3)

	// Rule 4:  IF e = LP OR  delta = LP THEN output = LN
	eR = getOutputError(table2, LP)
	dR = getOutputDelta(table2, LP)
	m4 := math.Max(eR, dR)
	o4 := -1.0
	mx = append(mx, m4)
	outputs = append(outputs, o4)

	// Centroid calculation
	numerator := 0.0
	denominator := 0.0
	for i := 0; i < len(mx); i++ {
		numerator += mx[i] * outputs[i]
		denominator += mx[i]
	}
	u := numerator / denominator

	fmt.Printf("u = %.2f", u)
}

func fuzzyfication(s string, f float64) float64 {
	table := []X{}

	table = append(table, X{Name: LP, Input: -1, Output: 0})
	table = append(table, X{Name: LP, Input: -0.75, Output: 0})
	table = append(table, X{Name: LP, Input: -0.5, Output: 0})
	table = append(table, X{Name: LP, Input: -0.25, Output: 0})
	table = append(table, X{Name: LP, Input: 0.0, Output: 0})
	table = append(table, X{Name: LP, Input: 0.25, Output: 0})
	table = append(table, X{Name: LP, Input: 0.5, Output: 0.3})
	table = append(table, X{Name: LP, Input: 0.75, Output: 0.7})
	table = append(table, X{Name: LP, Input: 1, Output: 1.0})

	table = append(table, X{Name: SP, Input: -1, Output: 0})
	table = append(table, X{Name: SP, Input: -0.75, Output: 0})
	table = append(table, X{Name: SP, Input: -0.5, Output: 0})
	table = append(table, X{Name: SP, Input: -0.25, Output: 0})
	table = append(table, X{Name: SP, Input: 0.0, Output: 0.3})
	table = append(table, X{Name: SP, Input: 0.25, Output: 0.7})
	table = append(table, X{Name: SP, Input: 0.5, Output: 1.0})
	table = append(table, X{Name: SP, Input: 0.75, Output: 0.7})
	table = append(table, X{Name: SP, Input: 1, Output: 0.3})

	table = append(table, X{Name: ZE, Input: -1, Output: 0})
	table = append(table, X{Name: ZE, Input: -0.75, Output: 0})
	table = append(table, X{Name: ZE, Input: -0.5, Output: 0.3})
	table = append(table, X{Name: ZE, Input: -0.25, Output: 0.7})
	table = append(table, X{Name: ZE, Input: 0.0, Output: 1.0})
	table = append(table, X{Name: ZE, Input: 0.25, Output: 0.7})
	table = append(table, X{Name: ZE, Input: 0.5, Output: 0.3})
	table = append(table, X{Name: ZE, Input: 0.75, Output: 0.0})
	table = append(table, X{Name: ZE, Input: 1, Output: 0.0})

	table = append(table, X{Name: SN, Input: -1, Output: 0.3})
	table = append(table, X{Name: SN, Input: -0.75, Output: 0.7})
	table = append(table, X{Name: SN, Input: -0.5, Output: 1.0})
	table = append(table, X{Name: SN, Input: -0.25, Output: 0.7})
	table = append(table, X{Name: SN, Input: 0.0, Output: 0.3})
	table = append(table, X{Name: SN, Input: 0.25, Output: 0.0})
	table = append(table, X{Name: SN, Input: 0.5, Output: 0.0})
	table = append(table, X{Name: SN, Input: 0.75, Output: 0.0})
	table = append(table, X{Name: SN, Input: 1, Output: 0.0})

	table = append(table, X{Name: LN, Input: -1, Output: 1.0})
	table = append(table, X{Name: LN, Input: -0.75, Output: 0.7})
	table = append(table, X{Name: LN, Input: -0.5, Output: 0.3})
	table = append(table, X{Name: LN, Input: -0.25, Output: 0.0})
	table = append(table, X{Name: LN, Input: 0.0, Output: 0.0})
	table = append(table, X{Name: LN, Input: 0.25, Output: 0.0})
	table = append(table, X{Name: LN, Input: 0.5, Output: 0.0})
	table = append(table, X{Name: LN, Input: 0.75, Output: 0.0})
	table = append(table, X{Name: LN, Input: 1, Output: 0.0})

	output := getOutput(table, s, f)

	return output
}

func getOutput(table []X, s string, f float64) float64 {
	found := false
	r := 0.0

	for i := 0; i < len(table); i++ {
		if table[i].Name == s && table[i].Input == f {
			found = true
			r = table[i].Output
			break
		}
	}
	if !found {
		fmt.Println("Error:: SOmething is wrong with the table")
		os.Exit(0)
	}
	return r
}

func getOutputError(table []X, s string) float64 {
	found := false
	r1 := 0.0

	for i := 0; i < len(table); i++ {
		if table[i].Name == s {
			found = true
			r1 = table[i].Input // error
			break
		}
	}
	if !found {
		fmt.Println("Error:: Something is wrong with the table")
		os.Exit(0)
	}
	return r1
}

func getOutputDelta(table []X, s string) float64 {
	found := false
	r1 := 0.0

	for i := 0; i < len(table); i++ {
		if table[i].Name == s {
			found = true
			r1 = table[i].Output // error
			break
		}
	}
	if !found {
		fmt.Println("Error:: Something is wrong with the table")
		os.Exit(0)
	}
	return r1
}

func deffuzification(x [9]float64) {

	//	SUM( I = 1 TO 4 OF ( mu(I) * output(I) ) ) / SUM( I = 1 TO 4 OF mu(I) )
}
