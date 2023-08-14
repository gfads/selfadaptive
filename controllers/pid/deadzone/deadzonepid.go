/*********************************************************************************
Author: Nelson S Rosa
Description: This program implements a PID controller with dead zone as defined
			in "Feedback Control for Computer Systems: Introducing Control Theory to
			Enterprise Programmers", Philipp K. Janert, 2014.
Date: 04/02/2023
*********************************************************************************/

package deadzonepid

import (
	"fmt"
	"main.go/controllers/def/info"
	"main.go/shared"

	"math"
	"os"
)

type Controller struct {
	Info info.Controller
}

func (c *Controller) Initialise(p ...float64) {

	if len(p) < 7 {
		fmt.Printf("Error: '%s' controller requires 7 info (min,max,kp,ki,kd,dead zone band) \n", shared.DeadZonePid)
		os.Exit(0)
	}

	c.Info.Direction = p[0]
	c.Info.Min = p[1]
	c.Info.Max = p[2]

	c.Info.Kp = p[3]
	c.Info.Ki = p[4]
	c.Info.Kd = p[5]
	c.Info.DeadZone = p[6]

	c.Info.Integrator = 0.0
	c.Info.PreviousError = 0.0
	c.Info.PreviousPreviousError = 0.0
	c.Info.SumPrevErrors = 0.0
	c.Info.Out = 0.0
	c.Info.PreviousDifferentiator = 0.0
}

func (c *Controller) Update(p ...float64) float64 {

	r := p[0] // goal
	y := p[1] // plant output

	// errors
	err := c.Info.Direction * (r - y)

	if math.Abs(err) > c.Info.DeadZone { // outside deadzone
		// Proportional
		proportional := c.Info.Kp * err

		// Integrator (David page 49)
		integrator := (c.Info.SumPrevErrors + err) * c.Info.Ki * shared.DeltaTime

		// Differentiator (David page 49)
		differentiator := c.Info.Kd * (err - c.Info.PreviousError) / shared.DeltaTime

		// pid output
		c.Info.Out = proportional + integrator + differentiator
	} else { // inside deadzone
		c.Info.Out = c.Info.PreviousOut // No action(?) TODO
	}

	if c.Info.Out > c.Info.Max {
		c.Info.Out = c.Info.Max
	} else if c.Info.Out < c.Info.Min {
		c.Info.Out = c.Info.Min
	}

	c.Info.PreviousOut = c.Info.Out
	c.Info.PreviousError = err
	c.Info.SumPrevErrors += err

	return c.Info.Out
}

func (c *Controller) SetGains(p ...float64) {
	c.Info.Kp = p[0]
	c.Info.Ki = p[1]
	c.Info.Kd = p[2]
}
