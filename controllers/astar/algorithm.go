package algorithm

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
		fmt.Printf("Error: '%s' controller requires 6 info (direction,min,max,kp,ki,kd) \n", shared.AsTAR)
		os.Exit(0)
	}

	c.Info.Min = p[0]
	c.Info.Max = p[1]
	c.Info.HysteresisBand = p[2]
}

func (c *Controller) Update(p ...float64) float64 {
	u := 0.0
	setpoint := p[0]
	y := p[1] // measured arrival rate

	if y < (setpoint - c.Info.HysteresisBand) { // The system is bellow the goal
		if y > c.Info.PreviousRate {
			u = c.Info.PreviousOut + 1
			//fmt.Printf("Below the goal (Accelerating) [%.4f]", c.Info.OptimumLevel-c.Info.HysteresisBand)
		} else {
			u = c.Info.PreviousOut * 2
			//fmt.Printf("Below the goal (Accelerating fast) [%.4f]", c.Info.OptimumLevel-c.Info.HysteresisBand)
		}
	} else if y > (setpoint + c.Info.HysteresisBand) { // The system is above the goal
		if y < c.Info.PreviousRate {
			u = c.Info.PreviousOut - 1
			//fmt.Printf("Above the goal (Reducing) [%.4f]", c.Info.OptimumLevel+c.Info.HysteresisBand)
		} else {
			u = c.Info.PreviousOut / 2
			//fmt.Printf("Above the goal (Reducing fast) [%.4f]", c.Info.OptimumLevel+c.Info.HysteresisBand)
		}
	} else { // The system is at Optimum state, no action required
		u = c.Info.PreviousOut
		//fmt.Printf("Optimum Level ")
	}

	// final check of rnew
	if u < c.Info.Min {
		u = c.Info.Min
	}
	if u > c.Info.Max {
		u = c.Info.Max
	}

	//fmt.Printf("[Rate=%.4f -> %.4f], [PC=%.4f -> %.4f]\n", c.Info.PreviousRate, y, c.Info.PreviousOut, u)

	c.Info.PreviousOut = u
	c.Info.PreviousRate = y

	return u
}

func (c *Controller) SetGains(p ...float64) {
}
