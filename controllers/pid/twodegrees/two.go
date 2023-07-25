/*********************************************************************************
Author: Nelson S Rosa
Description: This program implements a simple PID controller with setpoint weighting
as defined in Feedback Control for Computer Systems by Philipp K. Janert
Date: 25/07/2023
*********************************************************************************/
package twodegrees

import (
	"fmt"
	"main.go/controllers/def/info"
	"main.go/shared"
	"os"
)

const DeltaTime = 1 // see page 103
const Alfa = 1.0    // Alfa > 0
const Beta = 0.5    // Beta < 1

type Controller struct {
	Info info.Controller
}

func (c *Controller) Initialise(p ...float64) {

	if len(p) < 7 {
		fmt.Printf("Error: '%s' controller requires 7 info (direction,min,max,kp,ki,kd,beta) \n", shared.BasicPid)
		os.Exit(0)
	}

	c.Info.Direction = p[0]
	c.Info.Min = p[1]
	c.Info.Max = p[2]

	c.Info.Kp = p[3]
	c.Info.Ki = p[4]
	c.Info.Kd = p[5]
	c.Info.Beta = p[6]

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

	// Proportional
	proportional := c.Info.Kp * c.Info.Direction * (Alfa*r - y)

	// Integrator (page 49)
	c.Info.Integrator += DeltaTime * err
	integrator := c.Info.Integrator * c.Info.Ki

	// Differentiator (page 108)
	differentiator := c.Info.Kd * ((1-Beta)*r - y - c.Info.PreviousError) / DeltaTime

	// control law
	c.Info.Out = proportional + integrator + differentiator

	if c.Info.Out > c.Info.Max {
		c.Info.Out = c.Info.Max
	} else if c.Info.Out < c.Info.Min {
		c.Info.Out = c.Info.Min
	}

	c.Info.PreviousError = err
	c.Info.SumPrevErrors += err

	return c.Info.Out
}

func (c *Controller) SetGains(p ...float64) {
	c.Info.Kp = p[0]
	c.Info.Ki = p[1]
	c.Info.Kd = p[2]
}