/*********************************************************************************
Author: Nelson S Rosa
Description: This program defines the generic interface implemented by all controllers.
Date: 04/02/2023
*********************************************************************************/

package ops

import (
	"fmt"
	"main.go/controllers/def/info"
	gainscheduling "main.go/controllers/gain"
	onoffbasic "main.go/controllers/onoff/basic"
	deadzoneonff "main.go/controllers/onoff/deadzone"
	hysteresisonoff "main.go/controllers/onoff/hysteresis"
	basicpid "main.go/controllers/pid/basic"
	deadzonepid "main.go/controllers/pid/deadzone"
	errorsquarepid "main.go/controllers/pid/errorsquare"
	incrementalpid "main.go/controllers/pid/incremental"
	smoothingpid "main.go/controllers/pid/smoothing"
	"main.go/shared"
	"os"
)

type IController interface {
	Initialise(...float64)     // Initialise the controller
	Update(...float64) float64 // Update the controller output
	SetGains(...float64)       // Configure gains of PID controllers
}

// Create a controller of 'Type' (typeName) and configure its parameters //

func NewController(i info.Controller) IController {

	switch i.TypeName {
	case shared.BasicOnoff:
		c := onoffbasic.Controller{}
		c.Initialise(i.Direction, i.Min, i.Max)
		return &c
	case shared.DeadZoneOnoff:
		c := deadzoneonff.Controller{}
		c.Initialise(i.Direction, i.Min, i.Max, i.DeadZone)
		return &c
	case shared.HysteresisOnoff:
		c := hysteresisonoff.Controller{}
		c.Initialise(i.Direction, i.Min, i.Max, i.HysteresisBand)
		return &c
	case shared.BasicP:
		c := basicpid.Controller{}
		c.Initialise(i.Direction, i.Min, i.Max, i.Kp, 0.0, 0.0)
		return &c
	case shared.BasicPi:
		c := basicpid.Controller{}
		c.Initialise(i.Direction, i.Min, i.Max, i.Kp, i.Ki, 0.0)
		return &c
	case shared.BasicPid:
		c := basicpid.Controller{}
		c.Initialise(i.Direction, i.Min, i.Max, i.Kp, i.Ki, i.Kd)
		return &c
	case shared.SmoothingPid:
		c := smoothingpid.Controller{}
		c.Initialise(i.Direction, i.Min, i.Max, i.Kp, i.Ki, i.Kd)
		return &c
	case shared.IncrementalFormPid:
		c := incrementalpid.Controller{}
		c.Initialise(i.Direction, i.Min, i.Max, i.Kp, i.Ki, i.Kd)
		return &c
	case shared.DeadZonePid:
		c := deadzonepid.Controller{}
		c.Initialise(i.Direction, i.Min, i.Max, i.Kp, i.Ki, i.Kd, i.DeadZone)
		return &c
	case shared.ErrorSquarePid:
		c := errorsquarepid.Controller{}
		c.Initialise(i.Direction, i.Min, i.Max, i.Kp, i.Ki, i.Kd)
		return &c
	case shared.GainScheduling:
		c := gainscheduling.Controller{}
		c.Initialise(i.Direction, i.Min, i.Max, i.Kp, i.Ki, i.Kd)
		return &c
	default:
		fmt.Println("Error: Controller type ´", i.TypeName, "´ is unknown!")
		os.Exit(0)
	}

	return *new(IController)
}
