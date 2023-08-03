package parameters

import (
	"flag"
	"fmt"
	"main.go/shared"
	"os"
)

type ExecutionParameters struct {
	Tunning         *string
	ExecutionType   *string
	IsAdaptive      *bool
	ControllerType  *string
	MonitorInterval *int
	SetPoint        *float64
	Kp              *float64
	Ki              *float64
	Kd              *float64
	PrefetchCount   *int
	Min             *float64
	Max             *float64
	DeadZone        *float64
	HysteresisBand  *float64
	Direction       *float64
	GainTrigger     *float64
	Alfa            *float64
	Beta            *float64
}

func (e ExecutionParameters) Load() ExecutionParameters {
	p := ExecutionParameters{}

	p.ExecutionType = flag.String("execution-type", shared.StaticGoal, "execution-type is a string")
	p.IsAdaptive = flag.Bool("is-adaptive", false, "is-adaptive is a boolean")
	p.ControllerType = flag.String("controller-type", "OnOff", "controller-type is a string")
	p.MonitorInterval = flag.Int("monitor-interval", 1, "monitor-interval is an int (s)")
	p.SetPoint = flag.Float64("set-point", 3000.0, "set-point is a float (goal rate)")
	p.Kp = flag.Float64("kp", 1.0, "Kp is a float")
	p.Ki = flag.Float64("ki", 1.0, "Ki is a float")
	p.Kd = flag.Float64("kd", 1.0, "Kd is a float")
	p.PrefetchCount = flag.Int("prefetch-count", 1, "prefetch-count is an int")
	p.Min = flag.Float64("min", 0.0, "min is a float")
	p.Max = flag.Float64("max", 100.0, "max is a float")
	p.DeadZone = flag.Float64("dead-zone", 0.0, "dead-zone is a float")
	p.HysteresisBand = flag.Float64("hysteresis-band", 0.0, "hysteresis-band is a float")
	p.Direction = flag.Float64("direction", 1.0, "direction is a float")
	p.GainTrigger = flag.Float64("gain-trigger", 1.0, "gain trigger is a float")
	p.Alfa = flag.Float64("alfa", 1.0, "Alfa is a float (Setpoint Weighting)")
	p.Beta = flag.Float64("beta", 1.0, "Beta is a float (Setpoint Weighting / Two degrees of freedom)")
	p.Tunning = flag.String("tunning", "RootLocus", "tunning-type is a string")
	flag.Parse()

	return p
}

func (e ExecutionParameters) Validate(p ExecutionParameters) {
	if *p.Direction != 1.0 && *p.Direction != -1.0 {
		shared.ErrorHandler(shared.GetFunction(), "Direction invalid")
	}

	if *p.Tunning != shared.RootLocus && *p.Tunning != shared.Ziegler && *p.Tunning != shared.Cohen && *p.Tunning != shared.Amigo && *p.Tunning != shared.None {
		shared.ErrorHandler(shared.GetFunction(), "Tunning ´"+*p.Tunning+"´ is invalid")
	}
}

func (e ExecutionParameters) Show(p ExecutionParameters) {

	// validate execution type
	fmt.Println("************************************************")
	fmt.Printf("Execution Type  : %v\n", *p.ExecutionType)
	fmt.Printf("Is Adaptive?    : %v\n", *p.IsAdaptive)
	fmt.Printf("Tunning         : %v\n", *p.Tunning)
	fmt.Printf("Controller Type : %v\n", *p.ControllerType)
	fmt.Printf("Monitor Interval: %v\n", *p.MonitorInterval)
	fmt.Printf("Goal            : %.4f\n", *p.SetPoint)
	fmt.Printf("Prefetch Count  : %v\n", *p.PrefetchCount)
	fmt.Printf("Direction       : %.1f\n", *p.Direction)

	switch *p.ControllerType {
	case shared.AsTAR:
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
		fmt.Printf("Hysteresis Band : %.4f\n", *p.HysteresisBand)
	case shared.BasicOnoff:
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
	case shared.DeadZoneOnoff:
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
		fmt.Printf("Dead Zone       : %.4f\n", *p.DeadZone)
	case shared.HysteresisOnoff:
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
		fmt.Printf("Hystereis Band  : %.4f\n", *p.HysteresisBand)
	case shared.BasicP:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
	case shared.BasicPi:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
	case shared.BasicPid:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
	case shared.SmoothingPid:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
	case shared.IncrementalFormPid:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
	case shared.ErrorSquarePidFull:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
	case shared.ErrorSquarePidProportional:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
	case shared.DeadZonePid:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
		fmt.Printf("Dead Zone       : %.4f\n", *p.DeadZone)
	case shared.GainScheduling:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
		fmt.Printf("Gain Trigger    : %.4f\n", *p.GainTrigger)
	case shared.HPA:
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
		fmt.Printf("PC           : %v\n", *p.PrefetchCount)
	case shared.PIwithTwoDegreesOfFreedom:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
		fmt.Printf("Beta            : %.4f\n", *p.Beta)
	case shared.WindUp:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
	case shared.SetpointWeighting:
		fmt.Printf("Kp              : %.8f\n", *p.Kp)
		fmt.Printf("Ki              : %.8f\n", *p.Ki)
		fmt.Printf("Kd              : %.8f\n", *p.Kd)
		fmt.Printf("Min             : %.4f\n", *p.Min)
		fmt.Printf("Max             : %.4f\n", *p.Max)
		fmt.Printf("Alpha (Integral): %.4f\n", *p.Alfa)
		fmt.Printf("Beta (Derivative): %.4f\n", *p.Beta)

	default:
		fmt.Println(shared.GetFunction(), "Controller type ´", *p.ControllerType, "´ is invalid")
		os.Exit(0)
	}
	fmt.Println("************************************************")
}
