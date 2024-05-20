/*********************************************************************************
Author: Nelson S Rosa
Description: This program implements the On Off controller as defined in "Feedback
			Control for Computer Systems: Introducing Control Theory to Enterprise
			Programmers", Philipp K. Janert, 2014.

Date: 04/02/2023
*********************************************************************************/
package hpa

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

	if len(p) < 4 {
		fmt.Printf("Error: '%s controller requires 4 info (direction,min,max) \n", shared.HPA)
		os.Exit(0)
	}

	c.Info.TypeName = shared.HPA
	c.Info.Direction = p[0]
	c.Info.Min = p[1]
	c.Info.Max = p[2]
	c.Info.PC = p[3]
	//fmt.Printf("Direction=%v Min=%.4f Max=%.4f PC=%v \n", c.Info.Direction, c.Info.Min, c.Info.Max, c.Info.PC)
}

func (c *Controller) Update(p ...float64) float64 {
	u := 0.0

	s := p[0] // goal
	y := p[1] // plant output

	u = math.Round(c.Info.PC * s / y)

	//fmt.Printf("s/y=%.4f round(s/y)=%.4f\n", c.Info.PC*(y/s), c.Info.PC*(s/y))

	fmt.Printf("HPA:: Goal=%v y=%.4f PC=%.4f u=%.4f \n", s, y, c.Info.PC, u)

	// control law
	if u > c.Info.Max {
		u = c.Info.Max
	} else if u < c.Info.Min {
		u = c.Info.Min
	}

	c.Info.PC = u

	return u
}

func (c *Controller) SetGains(p ...float64) {
}
