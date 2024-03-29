/*********************************************************************************
Author: Nelson S Rosa
Description: This program implements the On Off controller as defined in "Feedback
			Control for Computer Systems: Introducing Control Theory to Enterprise
			Programmers", Philipp K. Janert, 2014.

Date: 04/02/2023
*********************************************************************************/

package onoffbasic

import (
	"fmt"
	"main.go/controllers/def/info"
	"main.go/shared"
	"os"
)

type Controller struct {
	Info info.Controller
}

func (c *Controller) Initialise(p ...float64) {

	if len(p) < 3 {
		fmt.Printf("Error: '%s controller requires 3 info (direction,min,max) \n", shared.BasicOnoff)
		os.Exit(0)
	}

	c.Info.TypeName = shared.BasicOnoff
	c.Info.Direction = p[0]
	c.Info.Min = p[1]
	c.Info.Max = p[2]
}

func (c *Controller) Update(p ...float64) float64 {
	u := 0.0

	s := p[0] // goal
	y := p[1] // plant output

	// error
	err := c.Info.Direction * (s - y)

	// control law
	if err >= 0 { // lower than the goal
		u = c.Info.Max
	} else { // higher than the goal
		u = c.Info.Min
	}
	return u
}

func (c *Controller) SetGains(p ...float64) {
}
