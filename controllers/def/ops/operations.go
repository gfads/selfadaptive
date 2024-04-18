/*********************************************************************************
Author: Nelson S Rosa
Description: This program defines the generic interface implemented by all controllers.
Date: 04/02/2023
*********************************************************************************/

package ops

import (
	"fmt"
	algorithm "main.go/controllers/astar"
	"main.go/controllers/def/parameters"
	"main.go/controllers/fuzzy"
	gainscheduling "main.go/controllers/gain"
	"main.go/controllers/hpa"
	onoffbasic "main.go/controllers/onoff/basic"
	deadzoneonff "main.go/controllers/onoff/deadzone"
	hysteresisonoff "main.go/controllers/onoff/hysteresis"
	basicpid "main.go/controllers/pid/basic"
	deadzonepid "main.go/controllers/pid/deadzone"
	"main.go/controllers/pid/errorsquarefull"
	"main.go/controllers/pid/errorsquareproportional"
	incrementalpid "main.go/controllers/pid/incremental"
	"main.go/controllers/pid/setpointweighting"
	smoothingpid "main.go/controllers/pid/smoothing"
	"main.go/controllers/pid/twodegrees"
	"main.go/controllers/pid/windup"
	"main.go/shared"
	"os"
)

type IController interface {
	Initialise(...float64)     // Initialise the controller
	Update(...float64) float64 // Update the controller output
	SetGains(...float64)       // Configure gains of PID controllers
}

// Create a controller of 'Type' (typeName) and configure its parameters //
func NewController(p parameters.ExecutionParameters) IController {
	switch *p.ControllerType {
	case shared.HPA:
		c := hpa.Controller{}
		c.Initialise(*p.Direction, *p.Min, *p.Max, float64(*p.PrefetchCount))
		return &c
	case shared.AsTAR:
		c := algorithm.Controller{}
		c.Initialise(*p.Min, *p.Max, *p.HysteresisBand)
		return &c
	case shared.BasicOnoff:
		c := onoffbasic.Controller{}
		c.Initialise(*p.Direction, *p.Min, *p.Max)
		return &c
	case shared.DeadZoneOnoff:
		c := deadzoneonff.Controller{}
		c.Initialise(*p.Direction, *p.Min, *p.Max, *p.DeadZone)
		return &c
	case shared.HysteresisOnoff:
		c := hysteresisonoff.Controller{}
		c.Initialise(*p.Direction, *p.Min, *p.Max, *p.HysteresisBand)
		return &c
	case shared.BasicP:
		c := basicpid.Controller{}
		c.Initialise(*p.Direction, *p.Min, *p.Max, *p.Kp, 0.0, 0.0)
		return &c
	case shared.BasicPi:
		c := basicpid.Controller{}
		c.Initialise(*p.Direction, *p.Min, *p.Max, *p.Kp, *p.Ki, 0.0)
		return &c
	case shared.BasicPid:
		c := basicpid.Controller{}
		c.Initialise(*p.Direction, *p.Min, *p.Max, *p.Kp, *p.Ki, *p.Kd)
		return &c
	case shared.SmoothingPid:
		c := smoothingpid.Controller{}
		c.Initialise(*p.Direction, *p.Min, *p.Max, *p.Kp, *p.Ki, *p.Kd)
		return &c
	case shared.IncrementalFormPid:
		c := incrementalpid.Controller{}
		c.Initialise(*p.Direction, *p.Min, *p.Max, *p.Kp, *p.Ki, *p.Kd)
		return &c
	case shared.DeadZonePid:
		c := deadzonepid.Controller{}
		c.Initialise(*p.Direction, *p.Min, *p.Max, *p.Kp, *p.Ki, *p.Kd, *p.DeadZone)
		return &c
	case shared.ErrorSquarePidFull:
		c := errorsquarefull.Controller{}
		c.Initialise(*p.Direction, *p.Min, *p.Max, *p.Kp, *p.Ki, *p.Kd)
		return &c
	case shared.ErrorSquarePidProportional:
		c := errorsquareproportional.Controller{}
		c.Initialise(*p.Direction, *p.Min, *p.Max, *p.Kp, *p.Ki, *p.Kd)
		return &c
	case shared.GainScheduling:
		c := gainscheduling.Controller{}
		c.Initialise(*p.Direction, *p.Min, *p.Max, *p.Kp, *p.Ki, *p.Kd)
		return &c
	case shared.PIwithTwoDegreesOfFreedom:
		c := twodegrees.Controller{}
		c.Initialise(*p.Direction, *p.Min, *p.Max, *p.Kp, *p.Ki, *p.Kd, *p.Beta)
		return &c
	case shared.WindUp:
		c := windup.Controller{}
		c.Initialise(*p.Direction, *p.Min, *p.Max, *p.Kp, *p.Ki, *p.Kd)
		return &c
	case shared.SetpointWeighting:
		c := setpointweighting.Controller{}
		c.Initialise(*p.Direction, *p.Min, *p.Max, *p.Kp, *p.Ki, *p.Kd, *p.Alfa, *p.Beta)
		return &c
	case shared.Fuzzy:
		c := fuzzy.Controller{}
		c.Initialise()
		return &c
	default:
		fmt.Println(shared.GetFunction(), "Error: Controller type ´", *p.ControllerType, "´ is unknown!")
		os.Exit(0)
	}
	return *new(IController)
}
