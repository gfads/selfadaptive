/*********************************************************************************
Author: Nelson S Rosa
Description: This program implements a simple PID controller as defined
			in "Feedback Control for Computer Systems: Introducing Control Theory to
			Enterprise Programmers", Philipp K. Janert, 2014.
Date: 04/02/2023
*********************************************************************************/

package basicpid

import (
	"fmt"
	"os"
	"selfadaptive/controllers/def/info"
	"selfadaptive/shared"
)

const DeltaTime = 1 // see page 103

type Controller struct {
	Info info.Controller
}

func (c *Controller) Initialise(p ...float64) {

	if len(p) < 5 {
		fmt.Printf("Error: '%s' controller requires 5 info (min,max,kp,ki,kd) \n", shared.BasicPid)
		os.Exit(0)
	}

	c.Info.Min = p[0]
	c.Info.Max = p[1]

	c.Info.Kp = p[2]
	c.Info.Ki = p[3]
	c.Info.Kd = p[4]

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
	err := r - y

	// Proportional
	proportional := c.Info.Kp * err

	// Integrator (page 49)
	c.Info.Integrator += DeltaTime * err
	integrator := c.Info.Integrator * c.Info.Ki

	// Differentiator (page 49)
	differentiator := c.Info.Kd * (err - c.Info.PreviousError) / DeltaTime

	// control law
	c.Info.Out = proportional + integrator + differentiator

	if c.Info.Out > c.Info.Max {
		c.Info.Out = c.Info.Max
	} else if c.Info.Out < c.Info.Min {
		c.Info.Out = c.Info.Min
	}

	c.Info.PreviousError = err
	c.Info.SumPrevErrors = c.Info.SumPrevErrors + err

	return c.Info.Out
}
