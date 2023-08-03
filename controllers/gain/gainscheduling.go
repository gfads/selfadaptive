/*********************************************************************************
Author: Nelson S Rosa
Description: This program implements the Gain Scheduling strategy (adaptive controller),
			as defined in "Feedback Control for Computer Systems: Introducing Control Theory to Enterprise
			Programmers", Philipp K. Janert, 2014.

Date: 04/02/2023
*********************************************************************************/

package gainscheduling

import (
	"fmt"
	"main.go/controllers/def/info"
	"main.go/shared"
	"os"
)

const DeltaTime = 1 // see page 103

type Controller struct {
	Info      info.Controller
	GainTable [2][3]float64
}

func (c *Controller) Initialise(p ...float64) {

	if len(p) < 6 {
		fmt.Printf("Error: '%s' controller requires 6 info (min,max,kp,ki,kd) \n", shared.GainScheduling)
		os.Exit(0)
	}

	// hard coded gain scheduling table P // TODO
	c.GainTable[0][0] = 0.1 // kp[0] // P
	c.GainTable[0][1] = 0.0 // ki[1]
	c.GainTable[0][2] = 0.0 // kd[2]

	//c.GainTable[1][0] = -9600 // kp[0] // PID
	//c.GainTable[1][1] = 0.5   // ki[1]
	//c.GainTable[1][2] = 0.01  // kd[2]

	// Root Locus
	c.GainTable[1][0] = -0.00144086 // kp[0] // PI
	c.GainTable[1][1] = 0.00248495  // ki[1]
	c.GainTable[1][2] = 0.00057789  // kd[2]

	kp := c.GainTable[0][0]
	ki := c.GainTable[0][1]
	kd := c.GainTable[0][2]

	c.Info.Direction = p[0]
	c.Info.Min = p[1]
	c.Info.Max = p[2]

	c.Info.Kp = kp
	c.Info.Ki = ki
	c.Info.Kd = kd

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

	// decide about the gain - based on the plant output
	//if y < c.Info.GainTrigger { // Full
	// decide about the gain - based on the error
	if y < r*0.1 { // Full
		c.Info.Kp = c.GainTable[0][0]
		c.Info.Ki = c.GainTable[0][1]
		c.Info.Kd = c.GainTable[0][2]
	} else { // Root Locus
		c.Info.Kp = c.GainTable[1][0]
		c.Info.Ki = c.GainTable[1][1]
		c.Info.Kd = c.GainTable[1][2]
	}

	// Proportional
	proportional := c.Info.Kp * err

	// Integrator (David page 49)
	c.Info.Integrator += DeltaTime * err
	integrator := c.Info.Integrator * c.Info.Ki

	// Differentiator (David page 49)
	differentiator := c.Info.Kd * (err - c.Info.PreviousError) / DeltaTime

	// pid output
	c.Info.Out = proportional + integrator + differentiator

	if c.Info.Out > c.Info.Max {
		c.Info.Out = c.Info.Max
	} else if c.Info.Out != 0 && c.Info.Out < c.Info.Min { // due to gain scheduling full = 0
		c.Info.Out = c.Info.Min
	}

	c.Info.PreviousError = err
	c.Info.SumPrevErrors = c.Info.SumPrevErrors + err

	return c.Info.Out
}

func (c *Controller) SetGains(p ...float64) {
	c.Info.Kp = p[0]
	c.Info.Ki = p[1]
	c.Info.Kd = p[2]
}
