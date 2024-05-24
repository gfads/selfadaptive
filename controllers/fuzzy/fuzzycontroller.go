package fuzzy

import (
	"fmt"
	"main.go/controllers/def/info"
	"main.go/controllers/fuzzy/deffuzification"
	"main.go/controllers/fuzzy/fuzzification"
	"main.go/shared"
	"os"
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

func (c *Controller) Update(p ...float64) float64 {
	goal := p[0]
	rate := p[1]

	e := goal - rate

	// 1. Fuzzification
	fuzzifiedSetError := fuzzyInput(e, fuzzification.GAUSSIAN)

	// 2. apply rules
	output := applyRules(fuzzifiedSetError)

	// 3. Deffuzification
	f := deffuzification.Centroid{}
	u := f.Deffuzify(output)

	return u
}

func (c *Controller) SetGains(p ...float64) {}

func applyRules(e map[string]float64) shared.OutputX {
	o := shared.OutputX{}

	// Rule 1:  IF e = EXTREMELYPOSITIVE THEN output = LARGEINCREASE
	o.Mx = append(o.Mx, e[shared.EXTREMELYPOSITIVE])
	o.Out = append(o.Out, getMaxOutput(shared.LARGEINCREASE))

	// Rule 2:  IF error LARGEPOSITIVE THEN output = MEDIUMINCREASE
	o.Mx = append(o.Mx, e[shared.LARGEPOSITIVE])
	o.Out = append(o.Out, getMaxOutput(shared.MEDIUMINCREASE))

	// Rule 3:  IF e = SMALLPOSITIVE THEN output = SMALLINCREASE
	o.Mx = append(o.Mx, e[shared.SMALLPOSITIVE]) // saida = +1 s
	o.Out = append(o.Out, getMaxOutput(shared.SMALLINCREASE))

	// Rule 4:  IF e = ZE THEN output = MAINTAIN
	o.Mx = append(o.Mx, e[shared.ZERO]) // saida = 0
	o.Out = append(o.Out, getMaxOutput(shared.MAINTAIN))

	// Rule 5:  IF e = SMALLNEGATIVE THEN output = SMALLDECREASE
	o.Mx = append(o.Mx, e[shared.SMALLNEGATIVE]) // saida = -1 s
	o.Out = append(o.Out, getMaxOutput(shared.SMALLDECREASE))

	// Rule 6:  IF e = LARGENEGATIVE THEN output = MEDIUMDECREASE
	o.Mx = append(o.Mx, e[shared.LARGENEGATIVE])
	o.Out = append(o.Out, getMaxOutput(shared.MEDIUMDECREASE))

	// Rule 7:  IF e = EXTREMELYNEGATIVE THEN output = LARGEDECREASE
	o.Mx = append(o.Mx, e[shared.EXTREMELYNEGATIVE])
	o.Out = append(o.Out, getMaxOutput(shared.LARGEDECREASE))

	//fmt.Printf("[%.2f %.2f %.2f %.2f %.2f %.2f %.2f]\n", o.Mx[0], o.Mx[1], o.Mx[2], o.Mx[3], o.Mx[4], o.Mx[5], o.Mx[6])
	//fmt.Printf("[%.2f %.2f %.2f %.2f %.2f %.2f %.2f]\n", o.Out[0], o.Out[1], o.Out[2], o.Out[3], o.Out[4], o.Out[5], o.Out[6])
	return o
}
func fuzzyInput(x float64, mf string) map[string]float64 {
	r := map[string]float64{}

	switch mf {
	case fuzzification.TRIANGULAR:
		f := fuzzification.Triangular{}
		r[shared.EXTREMELYPOSITIVE] = f.Fuzzify(x, 1250, 5000, 10000)
		r[shared.LARGEPOSITIVE] = f.Fuzzify(x, 500, 1250, 2000)          //500,750,1000
		r[shared.SMALLPOSITIVE] = f.Fuzzify(x, 0, 625, 1250)             // 0, 500,1000
		r[shared.ZERO] = f.Fuzzify(x, -500, 0, 500)                      // -500,0,500
		r[shared.SMALLNEGATIVE] = f.Fuzzify(x, -1250, -625, 0)           //-1000,-500,0
		r[shared.LARGENEGATIVE] = f.Fuzzify(x, -2000, -1250, -500)       // -1000,-750,-500
		r[shared.EXTREMELYNEGATIVE] = f.Fuzzify(x, -1250, -5000, -10000) // -1000,-750,-500
	case fuzzification.GAUSSIAN:
		f := fuzzification.Gaussian{}
		r[shared.EXTREMELYPOSITIVE] = f.Fuzzify(x, 3000.0, 0.01)
		r[shared.LARGEPOSITIVE] = f.Fuzzify(x, 1500.0, 0.01)      //500,750,1000
		r[shared.SMALLPOSITIVE] = f.Fuzzify(x, 500.0, 0.01)       // 0, 500,1000
		r[shared.ZERO] = f.Fuzzify(x, 0.0, 0.1)                   // -500,0,500
		r[shared.SMALLNEGATIVE] = f.Fuzzify(x, -500.0, 0.01)      //-1000,-500,0
		r[shared.LARGENEGATIVE] = f.Fuzzify(x, -1500.0, 0.01)     // -1000,-750,-500
		r[shared.EXTREMELYNEGATIVE] = f.Fuzzify(x, -3000.0, 0.01) // -1000,-750,-500
	case fuzzification.PI:
		f := fuzzification.Pi{}
		r[shared.EXTREMELYPOSITIVE] = f.Fuzzify(x, 1250, 2500, 5000, 10000)
		r[shared.LARGEPOSITIVE] = f.Fuzzify(x, 500, 250, 1750, 2000)            //500,750,1000
		r[shared.SMALLPOSITIVE] = f.Fuzzify(x, 0, 250, 1000, 1250)              // 0, 500,1000
		r[shared.ZERO] = f.Fuzzify(x, -500, -250, 250, 500)                     // -500,0,500
		r[shared.SMALLNEGATIVE] = f.Fuzzify(x, -1250, -1000, -250, 0)           //-1000,-500,0
		r[shared.LARGENEGATIVE] = f.Fuzzify(x, -2000, -1750, -250, -500)        // -1000,-750,-500
		r[shared.EXTREMELYNEGATIVE] = f.Fuzzify(x, -10000, -5000, -2500, -1250) // -1000,-750,-500
	default:
		fmt.Println("Error: Membership function invalid!")
		os.Exit(0)
	}

	/*
		fmt.Printf("Error = %.2f FuzzifiedError [%.2f %.2f %.2f %.2f %.2f %.2f %.2f]\n", x,
			r[shared.EXTREMELYNEGATIVE],
			r[shared.LARGENEGATIVE],
			r[shared.SMALLNEGATIVE],
			r[shared.ZERO],
			r[shared.SMALLPOSITIVE],
			r[shared.LARGEPOSITIVE],
			r[shared.EXTREMELYPOSITIVE])
	*/
	return r
}
func fuzzyOutput(n float64, mf string) map[string]float64 {
	r := map[string]float64{}

	switch mf {

	case fuzzification.GAUSSIAN:
		f := fuzzification.Gaussian{}
		r[shared.LARGEINCREASE] = f.Fuzzify(n, 3.0, 0.01)  // original = 2
		r[shared.MEDIUMINCREASE] = f.Fuzzify(n, 2.0, 0.01) // original = 2
		r[shared.SMALLINCREASE] = f.Fuzzify(n, 1.0, 0.01)  // original = 1
		r[shared.MAINTAIN] = f.Fuzzify(n, 0.0, 0.01)
		r[shared.SMALLDECREASE] = f.Fuzzify(n, -1.0, 0.01)  // original=-1
		r[shared.MEDIUMDECREASE] = f.Fuzzify(n, -2.0, 0.01) // original=-1
		r[shared.LARGEDECREASE] = f.Fuzzify(n, -3.0, 0.01)  // original= -2
	case fuzzification.TRIANGULAR:
		f := fuzzification.Triangular{}
		r[shared.LARGEINCREASE] = f.Fuzzify(n, 1.0, 2.0, 3.0)
		r[shared.SMALLINCREASE] = f.Fuzzify(n, 0.0, 1.0, 2.0)
		r[shared.MAINTAIN] = f.Fuzzify(n, 0.5, 0.0, -0.5)
		r[shared.SMALLDECREASE] = f.Fuzzify(n, -2.0, -1.0, 0.0)
		r[shared.LARGEDECREASE] = f.Fuzzify(n, -3.0, -2.0, -1.0)
	case fuzzification.TRAPEZOIDAL:
		f := fuzzification.Trapezoidal{}
		r[shared.LARGEINCREASE] = f.Fuzzify(n, 2.0, 2.5, 3.5, 4.0)
		r[shared.SMALLINCREASE] = f.Fuzzify(n, 1.0, 1.5, 2.5, 3.0)
		r[shared.MAINTAIN] = f.Fuzzify(n, -1.0, -0.5, 0.5, 1.0)
		r[shared.SMALLDECREASE] = f.Fuzzify(n, -3.0, -2.5, -1.5, -1.0)
		r[shared.LARGEDECREASE] = f.Fuzzify(n, -4.0, -3.5, -2.5, -3.0)
	default:
		fmt.Println("Error: Membership function invalid!")
		os.Exit(0)
	}
	return r
}
func getMaxOutput(s string) float64 {
	r := 0.0
	max := -10000.0

	for i := -3.0; i <= 3.0; i += 0.5 { // TODO
		v := fuzzyOutput(i, fuzzification.GAUSSIAN)
		if v[s] > max {
			max = v[s]
			r = i
		}
	}
	return r
}
