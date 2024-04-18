package parameters

import (
	"flag"
	"fmt"
	"main.go/shared"
	"os"
	"strconv"
)

type ExecutionParameters struct {
	OutputFile      *string
	Tunning         *string
	ExecutionType   *string
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

	p.ExecutionType = flag.String("execution-type", shared.Experiment, "execution-type is a string")
	p.ControllerType = flag.String("controller-type", shared.Fuzzy, "controller-type is a string")
	p.MonitorInterval = flag.Int("monitor-interval", 10, "monitor-interval is an int (s)")
	p.SetPoint = flag.Float64("set-point", 600.0, "set-point is a float (goal rate)")
	p.Kp = flag.Float64("kp", 1.0, "Kp is a float")
	p.Ki = flag.Float64("ki", 1.0, "Ki is a float")
	p.Kd = flag.Float64("kd", 1.0, "Kd is a float")
	p.PrefetchCount = flag.Int("prefetch-count", 1, "prefetch-count is an int")
	p.Min = flag.Float64("min", 1.0, "min is a float")
	p.Max = flag.Float64("max", 100.0, "max is a float")
	p.DeadZone = flag.Float64("dead-zone", 0.0, "dead-zone is a float")
	p.HysteresisBand = flag.Float64("hysteresis-band", 0.0, "hysteresis-band is a float")
	p.Direction = flag.Float64("direction", 1.0, "direction is a float")
	p.GainTrigger = flag.Float64("gain-trigger", 1.0, "gain trigger is a float")
	p.Alfa = flag.Float64("alfa", 1.0, "Alfa is a float (Setpoint Weighting)")
	p.Beta = flag.Float64("beta", 1.0, "Beta is a float (Setpoint Weighting / Two degrees of freedom)")
	p.Tunning = flag.String("tunning", "RootLocus", "tunning-type is a string")
	p.OutputFile = flag.String("output-file", "apague1.csv", "output-file is a string")
	//flag.Parse()

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

	r := "************************************************ \n" +
		"Output File     : " + *p.OutputFile + "\n" +
		"Execution Type  : " + *p.ExecutionType + "\n" +
		"Tunning         : " + *p.Tunning + "\n" +
		"Controller Type : " + *p.ControllerType + "\n" +
		"Monitor Interval: " + strconv.Itoa(*p.MonitorInterval) + "\n" +
		"Goal            : " + strconv.FormatFloat(*p.SetPoint, 'f', -1, 32) + "\n" +
		"Prefetch Count  : " + strconv.Itoa(*p.PrefetchCount) + "\n" +
		"Direction       : " + strconv.FormatFloat(*p.Direction, 'f', -1, 32) + "\n" +
		"Min             : " + strconv.FormatFloat(*p.Min, 'f', -1, 32) + "\n" +
		"Max             : " + strconv.FormatFloat(*p.Max, 'f', -1, 32) + "\n"

	switch *p.ControllerType {
	case shared.None:
	case shared.AsTAR:
		r += "Hysteresis Band : " + strconv.FormatFloat(*p.HysteresisBand, 'E', -1, 32)
	case shared.BasicOnoff:
	case shared.DeadZoneOnoff:
		r += "Dead Zone       : " + strconv.FormatFloat(*p.DeadZone, 'E', -1, 32) + "\n"
	case shared.HysteresisOnoff:
		r += "Hysteresis Band  : " + strconv.FormatFloat(*p.DeadZone, 'E', -1, 32) + "\n"
	case shared.BasicP:
		r += "Kp              : " + strconv.FormatFloat(*p.Kp, 'E', -1, 32) + "\n"
		r += "Ki              : " + strconv.FormatFloat(*p.Ki, 'E', -1, 32) + "\n"
		r += "Kd              : " + strconv.FormatFloat(*p.Kd, 'E', -1, 32) + "\n"
	case shared.BasicPi:
		r += "Kp              : " + strconv.FormatFloat(*p.Kp, 'E', -1, 32) + "\n"
		r += "Ki              : " + strconv.FormatFloat(*p.Ki, 'E', -1, 32) + "\n"
		r += "Kd              : " + strconv.FormatFloat(*p.Kd, 'E', -1, 32) + "\n"
	case shared.BasicPid:
		r += "Kp              : " + strconv.FormatFloat(*p.Kp, 'E', -1, 32) + "\n"
		r += "Ki              : " + strconv.FormatFloat(*p.Ki, 'E', -1, 32) + "\n"
		r += "Kd              : " + strconv.FormatFloat(*p.Kd, 'E', -1, 32) + "\n"
	case shared.SmoothingPid:
		r += "Kp              : " + strconv.FormatFloat(*p.Kp, 'E', -1, 32) + "\n"
		r += "Ki              : " + strconv.FormatFloat(*p.Ki, 'E', -1, 32) + "\n"
		r += "Kd              : " + strconv.FormatFloat(*p.Kd, 'E', -1, 32) + "\n"
	case shared.IncrementalFormPid:
		r += "Kp              : " + strconv.FormatFloat(*p.Kp, 'E', -1, 32) + "\n"
		r += "Ki              : " + strconv.FormatFloat(*p.Ki, 'E', -1, 32) + "\n"
		r += "Kd              : " + strconv.FormatFloat(*p.Kd, 'E', -1, 32) + "\n"
	case shared.ErrorSquarePidFull:
		r += "Kp              : " + strconv.FormatFloat(*p.Kp, 'E', -1, 32) + "\n"
		r += "Ki              : " + strconv.FormatFloat(*p.Ki, 'E', -1, 32) + "\n"
		r += "Kd              : " + strconv.FormatFloat(*p.Kd, 'E', -1, 32) + "\n"
	case shared.ErrorSquarePidProportional:
		r += "Kp              : " + strconv.FormatFloat(*p.Kp, 'E', -1, 32) + "\n"
		r += "Ki              : " + strconv.FormatFloat(*p.Ki, 'E', -1, 32) + "\n"
		r += "Kd              : " + strconv.FormatFloat(*p.Kd, 'E', -1, 32) + "\n"
	case shared.DeadZonePid:
		r += "Kp              : " + strconv.FormatFloat(*p.Kp, 'E', -1, 32) + "\n"
		r += "Ki              : " + strconv.FormatFloat(*p.Ki, 'E', -1, 32) + "\n"
		r += "Kd              : " + strconv.FormatFloat(*p.Kd, 'E', -1, 32) + "\n"
		r += "Dead Zone       : " + strconv.FormatFloat(*p.DeadZone, 'E', -1, 32) + "\n"
	case shared.GainScheduling:
		r += "Kp              : " + strconv.FormatFloat(*p.Kp, 'E', -1, 32) + "\n"
		r += "Ki              : " + strconv.FormatFloat(*p.Ki, 'E', -1, 32) + "\n"
		r += "Kd              : " + strconv.FormatFloat(*p.Kd, 'E', -1, 32) + "\n"
		r += "Gain Trigger    : " + strconv.FormatFloat(*p.GainTrigger, 'E', -1, 32) + "\n"
	case shared.HPA:
		r += "PC           : " + strconv.Itoa(*p.PrefetchCount) + "\n"
	case shared.PIwithTwoDegreesOfFreedom:
		r += "Kp              : " + strconv.FormatFloat(*p.Kp, 'E', -1, 32) + "\n"
		r += "Ki              : " + strconv.FormatFloat(*p.Ki, 'E', -1, 32) + "\n"
		r += "Kd              : " + strconv.FormatFloat(*p.Kd, 'E', -1, 32) + "\n"
		r += "Beta            : " + strconv.FormatFloat(*p.Beta, 'E', -1, 32) + "\n"
	case shared.WindUp:
		r += "Kp              : " + strconv.FormatFloat(*p.Kp, 'E', -1, 32) + "\n"
		r += "Ki              : " + strconv.FormatFloat(*p.Ki, 'E', -1, 32) + "\n"
		r += "Kd              : " + strconv.FormatFloat(*p.Kd, 'E', -1, 32) + "\n"
	case shared.SetpointWeighting:
		r += "Kp               : " + strconv.FormatFloat(*p.Kp, 'E', -1, 32) + "\n"
		r += "Ki               : " + strconv.FormatFloat(*p.Ki, 'E', -1, 32) + "\n"
		r += "Kd               : " + strconv.FormatFloat(*p.Kd, 'E', -1, 32) + "\n"
		r += "Alpha (Integral) : " + strconv.FormatFloat(*p.Alfa, 'E', -1, 32) + "\n"
		r += "Beta (Derivative): " + strconv.FormatFloat(*p.Beta, 'E', -1, 32) + "\n"
	case shared.Fuzzy:
	default:
		fmt.Println(shared.GetFunction(), "Controller type ´", *p.ControllerType, "´ is invalid")
		os.Exit(0)
	}
	r += "************************************************"
	fmt.Println(r)
}
