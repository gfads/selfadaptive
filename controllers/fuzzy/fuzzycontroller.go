package fuzzy

import (
	"main.go/controllers/def/info"
	"main.go/shared"
	"math"
)

type Controller struct {
	Info info.Controller
}

const EXTREMELYPOSITIVE = "EP"
const LARGEPOSITIVE = "LP" // Large Positive
const SMALLPOSITIVE = "SP" // Small Positive
const ZERO = "ZE"          // Zero
const SMALLNEGATIVE = "SN" // Smal Negative
const LARGENEGATIVE = "LN" // Large Negative
const EXTREMELYNEGATIVE = "EN"

const LARGEINCREASE = "LI"  // Large Positive
const SMALLINCREASE = "SI"  // Small Positive
const MAINTAIN = "MAINTAIN" // Zero
const SMALLDECREASE = "SD"  // Small Negative
const LARGEDECREASE = "LD"  // Large Negative

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

	// Rule 6:  IF e = EXTREMELYPOSITIVE THEN output = LARGEINCREASE
	eR = e[EXTREMELYPOSITIVE]
	//dR = d[LP]
	//m4 := math.Max(eR, dR)
	m6 := eR
	o6 := getMaxOutput(LARGEINCREASE)
	mx = append(mx, m6)
	output = append(output, o6)

	// Rule 7:  IF e = EXTREMELYNEGATIVE THEN output = LARGEDECREASE
	eR = e[EXTREMELYNEGATIVE]
	//dR = d[LP]
	//m4 := math.Max(eR, dR)
	m7 := eR
	o7 := getMaxOutput(LARGEDECREASE)
	mx = append(mx, m7)
	output = append(output, o7)

	return mx, output
}

func fuzzyficationInput(x float64) map[string]float64 {

	r := map[string]float64{}
	/*
		r[EXTREMELYPOSITIVE] = triangularMF(x, 1250, 5000, 10000)
		r[LARGEPOSITIVE] = triangularMF(x, 500, 1250, 2000)          //500,750,1000
		r[SMALLPOSITIVE] = triangularMF(x, 0, 625, 1250)             // 0, 500,1000
		r[ZERO] = triangularMF(x, -500, 0, 500)                      // -500,0,500
		r[SMALLNEGATIVE] = triangularMF(x, -1250, -625, 0)           //-1000,-500,0
		r[LARGENEGATIVE] = triangularMF(x, -2000, -1250, -500)       // -1000,-750,-500
		r[EXTREMELYNEGATIVE] = triangularMF(x, -1250, -5000, -10000) // -1000,-750,-500
	*/

	r[EXTREMELYPOSITIVE] = gaussianMF(x, 5000.0, 0.01)
	r[LARGEPOSITIVE] = gaussianMF(x, 1250.0, 0.01)      //500,750,1000
	r[SMALLPOSITIVE] = gaussianMF(x, 625.0, 0.01)       // 0, 500,1000
	r[ZERO] = gaussianMF(x, 0.0, 0.01)                  // -500,0,500
	r[SMALLNEGATIVE] = gaussianMF(x, -625.0, 0.01)      //-1000,-500,0
	r[LARGENEGATIVE] = gaussianMF(x, -1250.0, 0.01)     // -1000,-750,-500
	r[EXTREMELYNEGATIVE] = gaussianMF(x, -5000.0, 0.01) // -1000,-750,-500

	/*
		r[EXTREMELYPOSITIVE] = piMF(x, 1250, 2500, 5000, 10000)
		r[LARGEPOSITIVE] = piMF(x, 500, 250, 1750, 2000)            //500,750,1000
		r[SMALLPOSITIVE] = piMF(x, 0, 250, 1000, 1250)              // 0, 500,1000
		r[ZERO] = piMF(x, -500, -250, 250, 500)                     // -500,0,500
		r[SMALLNEGATIVE] = piMF(x, -1250, -1000, -250, 0)           //-1000,-500,0
		r[LARGENEGATIVE] = piMF(x, -2000, -1750, -250, -500)        // -1000,-750,-500
		r[EXTREMELYNEGATIVE] = piMF(x, -10000, -5000, -2500, -1250) // -1000,-750,-500
	*/
	//fmt.Println(r[LP], r[SP], r[ZE], r[SN], r[LN])
	return r
}

func (c *Controller) Update(p ...float64) float64 {
	goal := p[0]
	rate := p[1]

	// Fuzzification
	// fuzzyfication
	e := goal - rate

	fuzzifiedSetError := fuzzyficationInput(e)
	//fuzzifiedSetDelta := fuzzyficationInput(d)

	// apply rules
	mx, output := applyRules(fuzzifiedSetError)

	// Deffuzification
	u := centroidDeffuzification(mx, output)

	//fmt.Printf("Fuzzy Controller: %.2f\n", u)
	return u
}

func (c *Controller) SetGains(p ...float64) {}

func rampMF(x float64, a float64, b float64) float64 {
	if x <= a {
		return 0
	} else if x >= b {
		return 1
	} else {
		return (x - a) / (b - a)
	}
}

func gaussianMF(x float64, m float64, d float64) float64 {
	temp := -math.Pow(x-m, 2.0) / 2.0 * math.Pow(d, 2.0)
	r := math.Exp(temp)

	return r
}

func piMF(x float64, a float64, b float64, c float64, d float64) float64 {
	r := 0.0

	if x <= a {
		r = 0.0
	}
	if a <= x && x <= ((a+b)/2.0) {
		r = 2 * math.Pow((x-a)/(b-a), 2.0)
	}
	if b <= x && x <= c {
		r = 1
	}
	if c <= x && x <= ((c+d)/2.0) {
		r = 1 - 2*math.Pow((x-c)/(d-c), 2.0)
	}
	if (c+d)/2 <= x && x <= d {
		r = 2 * math.Pow((x-d)/(d-c), 2.0)
	}
	if x >= d {
		r = 0
	}

	return r
}

func triangularMF(x float64, a float64, b float64, c float64) float64 {
	return math.Max(0, math.Min((x-a)/(b-a), (c-x)/(c-b)))
}

// Trapezoidal membership function
func trapezoidalMF(x float64, a float64, b float64, c float64, d float64) float64 {
	min1 := math.Min(x-a/b-a, 1)
	r := math.Max(math.Min(min1, d-x/d-c), 0)

	return r
}

func fuzzyficationOutput(n float64) map[string]float64 {

	r := map[string]float64{}

	r[LARGEINCREASE] = gaussianMF(n, 2.0, 0.01)
	r[SMALLINCREASE] = gaussianMF(n, 1.0, 0.01)
	r[MAINTAIN] = gaussianMF(n, 0.0, 0.01)
	r[SMALLDECREASE] = gaussianMF(n, -1.0, 0.01)
	r[LARGEDECREASE] = gaussianMF(n, -2.0, 0.01)

	/*
		r[LARGEINCREASE] = triangularMF(n, 2.0, 3.0, 4.0)
		r[SMALLINCREASE] = triangularMF(n, 1.0, 2.0, 3.0)
		r[MAINTAIN] = triangularMF(n, 1.0, 0.0, -1.0)
		r[SMALLDECREASE] = triangularMF(n, -1.0, -2.0, -3.0)
		r[LARGEDECREASE] = triangularMF(n, -2.0, -3.0, -4.0)
	*/
	/*
		r[LARGEINCREASE] = trapezoidalMF(n, 2.0, 2.5, 3.5, 4.0)
		r[SMALLINCREASE] = trapezoidalMF(n, 1.0, 1.5, 2.5, 3.0)
		r[MAINTAIN] = trapezoidalMF(n, -1.0, -0.5, 0.5, 1.0)
		r[SMALLDECREASE] = trapezoidalMF(n, -3.0, -2.5, -1.5, -1.0)
		r[LARGEDECREASE] = trapezoidalMF(n, -4.0, -3.5, -2.5, -3.0)
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
