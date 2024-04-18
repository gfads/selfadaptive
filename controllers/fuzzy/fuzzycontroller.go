package fuzzy

import (
	"main.go/controllers/def/info"
	"main.go/shared"
	"math"
)

type Controller struct {
	Info info.Controller
}

func (c *Controller) Initialise(p ...float64) {

	/*
		if len(p) < 4 {
			fmt.Printf("Error: '%s controller requires 4 info (input.center, input.width, output.center, output.width) \n", shared.Fuzzy)
			os.Exit(0)
		}*/

	c.Info.TypeName = shared.Fuzzy
	//c.Info.InputSet.Center = p[0]
	//c.Info.InputSet.Width = p[1]
	//c.Info.OutputSet.Center = p[2]
	//c.Info.OutputSet.Width = p[3]
}

// Define membership functions
func triangularMF(x float64, a float64, b float64, c float64) float64 {
	return math.Max(0, math.Min((math.Min((x-a)/(b-a), (c-x)/(c-b))), 1))
}

func maxMembershipDefuzzification(heaterIncrease float64, heaterMaintain float64, heaterDecrease float64) float64 {
	max := math.Max(heaterIncrease, math.Max(heaterMaintain, heaterDecrease))

	switch max {
	case heaterIncrease:
		return 10
	case heaterMaintain:
		return 5
	default:
		return 0
	}
}

func weightedAverageDefuzzification(heaterIncrease float64, heaterMaintain float64, heaterDecrease float64) float64 {
	weightedSum := (heaterIncrease*10 + heaterMaintain*5 + heaterDecrease*0)
	sum := heaterIncrease + heaterMaintain + heaterDecrease

	if sum == 0 {
		return 0
	}

	return weightedSum / sum
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

// Fuzzification: Fuzzy input mapping
func fuzzification(goal, rate float64) (float64, float64, float64, float64, float64, float64, float64) {

	bl := triangularMF(rate, goal*0.0, goal*0.25, 0.50*goal)
	ml := triangularMF(rate, goal*0.25, goal*0.50, 0.75*goal)
	sl := triangularMF(rate, goal*0.50, goal*0.75, goal)
	moderate := triangularMF(rate, goal*0.75, goal, goal*1.25)
	sh := triangularMF(rate, goal, goal*1.25, goal*1.50)
	mh := triangularMF(rate, goal*1.25, goal*1.50, goal*1.75)
	bh := triangularMF(rate, goal*1.50, goal*1.75, goal*2.0)

	/*bl := rampMF(rate, goal*0.10, 0.50*goal)
	ml := rampMF(rate, goal*0.25, 0.75*goal)
	sl := rampMF(rate, goal*0.50, goal)
	moderate := rampMF(rate, goal*0.75, goal*1.25)
	sh := rampMF(rate, goal, goal*1.50)
	mh := rampMF(rate, goal*1.25, goal*1.75)
	bh := rampMF(rate, goal*1.50, goal*2.0)
	*/
	/*
		bl := piMF(rate, goal*0.10, goal*0.20, goal*0.30, 0.40*goal)
		ml := piMF(rate, goal*0.30, goal*0.40, goal*0.50, 0.60*goal)
		sl := piMF(rate, goal*0.50, goal*0.60, goal*0.70, goal*0.80)
		moderate := piMF(rate, goal*0.70, goal*0.80, goal*0.90, goal)
		sh := piMF(rate, goal*0.90, goal, goal*1.10, goal*1.20)
		mh := piMF(rate, goal*1.10, goal*1.20, goal*1.30, goal*1.40)
		bh := piMF(rate, goal*1.30, goal*1.40, goal*1.50, goal*1.60)
	*/
	/*
		bl := trapezoidalMF(rate, goal*0.10, goal*0.25, goal*0.35, 0.5*goal)
		ml := trapezoidalMF(rate, goal*0.25, goal*0.50, goal*0.60, 0.75*goal)
		sl := trapezoidalMF(rate, goal*0.50, goal*0.75, goal*0.85, goal)
		moderate := trapezoidalMF(rate, goal*0.75, goal, goal*1.10, goal*1.25)
		sh := trapezoidalMF(rate, goal, goal*1.25, goal*1.35, goal*1.50)
		mh := trapezoidalMF(rate, goal*1.25, goal*1.50, goal*1.60, goal*1.75)
		bh := trapezoidalMF(rate, goal*1.50, goal*1.75, goal*1.85, goal*2.0)
	*/
	return bl, ml, sl, moderate, sh, mh, bh
}

func centroidDefuzzification(pcbl, pcml, pcsl, pcmoderate, pcsh, pcmh, pcbh float64) float64 {
	r := 0.0
	//numerator := pcbl*4.0 + pcml*3.0 + pcsl*2.0 + pcmoderate + pcsh/2.0 + pcmh/4.0 + pcbh/8.0 //initial
	numerator := pcbl*4 + pcml*3 + pcsl*2 + pcmoderate + pcsh/2.0 + pcmh/4.0 + pcbh/6.0
	denominator := pcbl + pcml + pcsl + pcmoderate + pcsh + pcmh + pcbh

	//fmt.Println("Deffuzification: ", numerator, denominator)
	if denominator == 0 {
		r = 1.0 // minumum PC
	} else {
		r = numerator / denominator
	}

	if math.Round(r) == 0 {
		r = 1
	}

	return math.Round(r)
}

func fuzzyRules(fr1 float64, fr2 float64, fr3 float64) (float64, float64, float64) {
	fpc1 := fr1
	fpc2 := fr2
	fpc3 := fr3

	return fpc1, fpc2, fpc3
}

func (c *Controller) Update(p ...float64) float64 {
	goal := p[0]
	rate := p[1]

	// Fuzzification
	bl, ml, sl, moderate, sh, mh, bh := fuzzification(goal, rate)

	//fmt.Printf("%.2f [%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f] \n", rate, bl, ml, sl, moderate, sh, mh, bh)

	// Fuzzy Rules
	pcbl := bl
	pcml := ml
	pcsl := sl
	pcmoderate := moderate
	pcsh := sh
	pcmh := mh
	pcbh := bh

	//outputFuzzy := fuzzyRules(inputFuzzy)

	// Defuzzification
	u := centroidDefuzzification(pcbl, pcml, pcsl, pcmoderate, pcsh, pcmh, pcbh)
	return u
}

func (c *Controller) SetGains(p ...float64) {
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
