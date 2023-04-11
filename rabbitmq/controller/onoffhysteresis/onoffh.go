package onoffhysteresis

type OnOffwithHysteresisController struct {
	Minimum        float64
	Maximum        float64
	HysteresisBand float64
	PreviousOut    float64
}

func (c *OnOffwithHysteresisController) Initialise(t string, p ...float64) {

	// TODO remove t
	c.Minimum = p[0]
	c.Maximum = p[1]
	c.HysteresisBand = p[2]
	c.PreviousOut = 0.0
}

func (c *OnOffwithHysteresisController) Reconfigure(p ...float64) {
	// TODO
}

func (c *OnOffwithHysteresisController) Update(p ...float64) float64 {
	direction := -1.0 // TODO
	u := 0.0

	s := p[0] // goal
	y := p[1] // arrival rate

	// error
	err := direction * (s - y)

	if err > -c.HysteresisBand/2.0 && err < c.HysteresisBand/2.0 {
		u = c.PreviousOut
	}
	if err >= c.HysteresisBand/2.0 {
		u = c.Maximum
	}
	if err <= -c.HysteresisBand/2.0 {
		u = c.Minimum
	}

	if u < c.Minimum {
		u = c.Minimum
	}
	if u > c.Maximum {
		u = c.Maximum
	}
	c.PreviousOut = u

	return u
}

func (c *OnOffwithHysteresisController) SetGains(kp, ki, kd float64) {
	// TODO remove
}

func (c *OnOffwithHysteresisController) SetKP(kp float64) {
	// TODO remove
}

func (c *OnOffwithHysteresisController) GetKP() float64 {
	// TODO remove
	return 0.0
}
