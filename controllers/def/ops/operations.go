/*********************************************************************************
Author: Nelson S Rosa
Description: This program defines the generic interface implemented by all controllers.
Date: 04/02/2023
*********************************************************************************/

package ops

import (
	"fmt"
	"os"
	"selfadaptive/controllers/def/info"
	gainscheduling "selfadaptive/controllers/gain"
	onoffbasic "selfadaptive/controllers/onoff/basic"
	deadzoneonff "selfadaptive/controllers/onoff/deadzone"
	hysteresisonoff "selfadaptive/controllers/onoff/hysteresis"
	basicpid "selfadaptive/controllers/pid/basic"
	deadzonepid "selfadaptive/controllers/pid/deadzone"
	errorsquarepid "selfadaptive/controllers/pid/errorsquare"
	incrementalpid "selfadaptive/controllers/pid/incremental"
	smoothingpid "selfadaptive/controllers/pid/smoothing"
	"selfadaptive/shared"
)

type IController interface {
	Initialise(...float64)     // Initialise the controller
	Update(...float64) float64 // Update the controller output
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
