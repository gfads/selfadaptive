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

	if len(p) < 4 {
		fmt.Printf("Error: '%s' controller requires 6 info (direction,min,max,kp,ki,kd) \n", shared.AsTAR)
		os.Exit(0)
	}

	c.Info.Min = p[0]
	c.Info.Max = p[1]

	c.Info.OptimumLevel = p[3]
	c.Info.ShutoffLevel = p[4]
}

func (a Controller) Update(vnew float64, vold float64, rold float64) float64 {
	rnew := 0.0
	//getnew := 0.0

	if vnew < a.Info.ShutoffLevel { // The system is in Shut-off Voltage state, task is stopped
		rnew = 0.0
		//fmt.Println("Shut-off voltage state", vnew, vold, rold, rnew)
	} else if vnew < (a.Info.OptimumLevel - a.Info.HysteresisBand) { // The system is in Low-voltage state, apply AIMD
		if vnew > vold {
			rnew = rold + 1
			//fmt.Println("Low-voltage voltage state (Accelerating)", vnew, vold, rold, rnew)
		} else {
			rnew = rold / 2
			//fmt.Println("Low-voltage voltage state (Reducing)", vnew, vold, rold, rnew)
		}
	} else if vnew > (a.Info.OptimumLevel + a.Info.HysteresisBand) { // The system is in High Voltage state, apply MIAD
		if vnew < vold {
			rnew = rold - 1
			//fmt.Println("High-voltage state", vnew, vold, rold, rnew)
		} else {
			rnew = rold * 2
			//fmt.Println("High-voltage state", vnew, vold, rold, rnew)
		}
	} else { // The system is at Optimum Voltage state, take no action
		rnew = rold
		//fmt.Println("Optimum voltage state", vnew, vold, rold, rnew)
	}

	// final check of rnew
	if rnew < a.Info.Min {
		rnew = a.Info.Min
	}
	if rnew > a.Info.Max {
		rnew = a.Info.Max
	}

	//if rnew != 0.0 {
	//	getnew = 1.0 / float64(rnew)
	//} else {
	//	getnew = 0.0
	//}

	return rnew
}
