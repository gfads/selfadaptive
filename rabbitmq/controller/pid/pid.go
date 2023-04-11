package pid

import (
	"fmt"
	"math"
	"os"
	"selfadaptive/rabbitmq/controller/info"
	"selfadaptive/shared"
)

const DELTA_TIME = 1 // see page 103
const ALPHA = 0.1    // alpha variates from 0 to 1 (see page 104)

type PIDController struct {
	Info info.InfoController
}

func (c *PIDController) Initialise(pidType string, p ...float64) {

	kp := p[0]
	ki := p[1]
	kd := p[2]

	limMin := p[3]
	limMax := p[4]

	if pidType == shared.DeadZonePid {
		c.Info.DeadZone = p[5]
	}

	c.Info.LimMin = limMin
	c.Info.LimMax = limMax

	c.Info.Kp = kp
	c.Info.Ki = ki
	c.Info.Kd = kd

	c.Info.Integrator = 0.0

	c.Info.PreviousError = 0.0
	c.Info.PreviousPreviousError = 0.0
	c.Info.SumPreviousErrors = 0.0
	c.Info.Out = 0.0
	c.Info.PreviousDifferentiator = 0.0
	c.Info.PIDType = pidType
}

func (c *PIDController) Reconfigure(p ...float64) {
	// TODO
}

func (c *PIDController) Update(p ...float64) float64 {
	r := 0.0

	switch c.Info.PIDType {
	case shared.BasicPid:
		r = c.UpdateBasic(p)
	case shared.SmoothingPid:
		r = c.UpdateSmoothingDerivative(p)
	case shared.IncrementalFormPid:
		r = c.UpdateIncrementalForm(p)
	case shared.ErrorSquarePid:
		r = c.UpdateErrorSquare(p) // Page 109
	case shared.DeadZonePid:
		r = c.UpdateDeadZonePID(p) // Page 110
	case shared.NonePid: // TODO
		r = c.UpdateBasic(p)
	default:
		fmt.Println("PID Type ´", c.Info.PIDType, "´ is invalid!!!")
		os.Exit(0)
	}
	return r
}

func (c *PIDController) UpdateBasic(p []float64) float64 {

	r := p[0] // goal
	y := p[1] // arrival rate

	// errors
	err := r - y

	// Proportional
	proportional := c.Info.Kp * err

	// Integrator (David page 49)
	c.Info.Integrator += DELTA_TIME * err
	integrator := c.Info.Integrator * c.Info.Ki

	// Differentiator (David page 49)
	differentiator := c.Info.Kd * (err - c.Info.PreviousError) / DELTA_TIME

	// pid output
	c.Info.Out = proportional + integrator + differentiator

	if c.Info.Out > c.Info.LimMax {
		c.Info.Out = c.Info.LimMax
	} else if c.Info.Out < c.Info.LimMin {
		c.Info.Out = c.Info.LimMin
	}

	c.Info.PreviousError = err
	c.Info.SumPreviousErrors = c.Info.SumPreviousErrors + err

	return c.Info.Out
}

func (c *PIDController) UpdateDeadZonePID(p []float64) float64 {

	r := p[0] // goal
	y := p[1] // arrival rate

	// errors
	err := r - y

	if math.Abs(err) < c.Info.DeadZone/2 {
		c.Info.Out = c.Info.Out // no change
	} else {
		// Proportional
		proportional := c.Info.Kp * err

		// Integrator (David page 49)
		c.Info.Integrator += DELTA_TIME * err
		integrator := c.Info.Integrator * c.Info.Ki

		// Differentiator (David page 49)
		differentiator := c.Info.Kd * (err - c.Info.PreviousError) / DELTA_TIME

		// pid output
		c.Info.Out = proportional + integrator + differentiator
	}

	if c.Info.Out > c.Info.LimMax {
		c.Info.Out = c.Info.LimMax
	} else if c.Info.Out < c.Info.LimMin {
		c.Info.Out = c.Info.LimMin
	}

	c.Info.PreviousError = err
	c.Info.SumPreviousErrors = c.Info.SumPreviousErrors + err

	return c.Info.Out
}

func (c *PIDController) UpdateErrorSquare(p []float64) float64 {

	r := p[0] // goal
	y := p[1] // arrival rate

	// errors
	err := r - y

	// Proportional
	proportional := c.Info.Kp * err

	// Integrator (David page 49)
	c.Info.Integrator += DELTA_TIME * err
	integrator := c.Info.Integrator * c.Info.Ki

	// Differentiator (David page 49)
	differentiator := c.Info.Kd * (err - c.Info.PreviousError) / DELTA_TIME

	// pid output Page 109
	c.Info.Out = math.Abs(err) * (proportional + integrator + differentiator)

	if c.Info.Out > c.Info.LimMax {
		c.Info.Out = c.Info.LimMax
	} else if c.Info.Out < c.Info.LimMin {
		c.Info.Out = c.Info.LimMin
	}

	c.Info.PreviousError = err
	c.Info.SumPreviousErrors = c.Info.SumPreviousErrors + err

	return c.Info.Out
}

func (c *PIDController) UpdateIncrementalForm(p []float64) float64 {

	r := p[0] // goal
	y := p[1] // arrival rate

	// errors
	err := r - y

	// Proportional // page 106
	//proportional := c.Info.Kp * err

	// Integrator // page 106
	c.Info.Integrator += DELTA_TIME * err
	integrator := c.Info.Integrator * c.Info.Ki

	// Differentiator // page 106
	//differentiator := c.Info.Kd * (err - c.Info.PreviousError) / DELTA_TIME

	// Delta of the new PC
	deltaU := c.Info.Kp*(err-c.Info.PreviousError) + c.Info.Ki*err*DELTA_TIME + c.Info.Kd*(err-2*c.Info.PreviousError+c.Info.PreviousPreviousError)/DELTA_TIME

	// pid output  TODO
	c.Info.Out = integrator + deltaU // see page 106 why add an integrator

	if c.Info.Out > c.Info.LimMax {
		c.Info.Out = c.Info.LimMax
	} else if c.Info.Out < c.Info.LimMin {
		c.Info.Out = c.Info.LimMin
	}

	c.Info.PreviousPreviousError = c.Info.PreviousError
	c.Info.PreviousError = err
	c.Info.SumPreviousErrors = c.Info.SumPreviousErrors + err

	return c.Info.Out
}

func (c *PIDController) UpdateSmoothingDerivative(p []float64) float64 {

	r := p[0] // goal
	y := p[1] // arrival rate

	// errors
	err := r - y

	// Proportional
	proportional := c.Info.Kp * err

	// Integrator (David page 49)
	c.Info.Integrator += DELTA_TIME * err
	integrator := c.Info.Integrator * c.Info.Ki

	// smoothing the derivative term (page 104)
	differentiator := ALPHA*(err-c.Info.PreviousError)/DELTA_TIME + (1-ALPHA)*c.Info.PreviousDifferentiator
	c.Info.PreviousDifferentiator = differentiator

	// pid output
	c.Info.Out = proportional + integrator + differentiator

	//fmt.Println("Kp=", c.Info.Kp, "Ki=", c.Info.Ki, "Kd=", c.Info.Kd)
	//fmt.Println("PID:: Goal=", r, "Rate=", y, "Error=", err, "[", proportional, integrator, differentiator, "]", c.Info.Out)

	if c.Info.Out > c.Info.LimMax {
		c.Info.Out = c.Info.LimMax
	} else if c.Info.Out < c.Info.LimMin {
		c.Info.Out = c.Info.LimMin
	}

	c.Info.PreviousError = err
	c.Info.SumPreviousErrors = c.Info.SumPreviousErrors + err

	return c.Info.Out
}

func (c *PIDController) SetGains(kp, ki, kd float64) {
	c.Info.Kp = kp
	c.Info.Ki = ki
	c.Info.Kd = kd
}

func (c *PIDController) SetKP(kp float64) {
	c.Info.Kp = kp
}

func (c *PIDController) GetKP() float64 {
	return c.Info.Kp
}
