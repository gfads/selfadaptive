package onoffdeadzone

type OnOffDeadZoneController struct {
	Minimum      float64
	Maximum      float64
	DeadZoneBand float64
}

func (c *OnOffDeadZoneController) Initialise(t string, p ...float64) {

	// TODO remove t
	c.Minimum = p[0]
	c.Maximum = p[1]
	c.DeadZoneBand = p[2]
}

func (c *OnOffDeadZoneController) Reconfigure(p ...float64) {
	// TODO
}

func (c *OnOffDeadZoneController) Update(p ...float64) float64 {

	direction := -1.0 // TODO
	u := 0.0

	s := p[0] // goal
	y := p[1] // arrival rate

	// error
	err := direction * (s - y)

	if err > -c.DeadZoneBand/2.0 && err < c.DeadZoneBand/2.0 {
		u = 0.0 // no action
	}
	if err >= c.DeadZoneBand/2.0 {
		u = c.Maximum
	}
	if err <= -c.DeadZoneBand/2 {
		u = c.Minimum
	}

	if u < c.Minimum {
		u = c.Minimum
	}
	if u > c.Maximum {
		u = c.Maximum
	}

	return u
}

func (c *OnOffDeadZoneController) SetGains(kp, ki, kd float64) {
	// TODO remove
}

func (c *OnOffDeadZoneController) SetKP(kp float64) {
	// TODO remove
}

func (c *OnOffDeadZoneController) GetKP() float64 {
	// TODO remove
	return 0.0
}
