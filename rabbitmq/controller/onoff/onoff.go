package onoff

type OnOffController struct {
	Minimum float64
	Maximum float64
}

func (c *OnOffController) Initialise(t string, p ...float64) {

	// TODO remove t
	c.Minimum = p[0]
	c.Maximum = p[1]
}

func (c *OnOffController) Reconfigure(p ...float64) {
	// TODO
}

func (c *OnOffController) Update(p ...float64) float64 {

	direction := -1.0 // TODO
	u := 0.0

	s := p[0] // goal
	y := p[1] // arrival rate

	// error
	err := direction * (s - y)

	if err >= 0 {
		u = c.Maximum
	} else {
		u = c.Minimum
	}
	return u
}

func (c *OnOffController) SetGains(kp, ki, kd float64) {
	// TODO remove
}

func (c *OnOffController) SetKP(kp float64) {
	// TODO remove
}

func (c *OnOffController) GetKP() float64 {
	// TODO remove
	return 0.0
}
