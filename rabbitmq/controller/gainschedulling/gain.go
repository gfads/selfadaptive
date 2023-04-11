package gain

import (
	"fmt"
	"math"
	"os"
	"selfadaptive/rabbitmq/controller/info"
	"selfadaptive/shared"
)

const DeltaTime = 1 // see page 103
const ALPHA = 0.1   // alpha variates from 0 to 1 (see page 104)

type GainPIDController struct {
	Info      info.InfoController
	GainTable [2][3]float64
}

func (c *GainPIDController) Initialise(pidType string, p ...float64) {

	// TODO
	c.GainTable[0][0] = -9600 // kp[0] // P
	c.GainTable[0][1] = 0.0   // ki[1]
	c.GainTable[0][2] = 0.0   // kd[2]

	//c.GainTable[1][0] = -9600 // kp[0] // PID
	//c.GainTable[1][1] = 0.5   // ki[1]
	//c.GainTable[1][2] = 0.01  // kd[2]

	c.GainTable[1][0] = -9600 // kp[0] // PI
	c.GainTable[1][1] = 0.5   // ki[1]
	c.GainTable[1][2] = 0.0   // kd[2]

	kp := c.GainTable[0][0]
	ki := c.GainTable[0][1]
	kd := c.GainTable[0][2]

	limMin := p[3]
	limMax := p[4]

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

func (c *GainPIDController) Reconfigure(p ...float64) {
	// TODO
	fmt.Println(p)
}

func (c *GainPIDController) Update(p ...float64) float64 {
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
	case shared.NonePid: // TODO
		r = c.UpdateBasic(p)
	default:
		fmt.Println("PID Type ´", c.Info.PIDType, "´ is invalid!!!")
		os.Exit(0)
	}
	return r
}

func (c *GainPIDController) UpdateBasic(p []float64) float64 {

	r := p[0] // goal
	y := p[1] // measured

	// errors
	err := r - y

	// decide about the gain -- based on the capacitor energy level
	if y < shared.OV { // gain scheduling 1
		//fmt.Println("************** HERE ****************")
		c.Info.Kp = c.GainTable[0][0]
		c.Info.Ki = c.GainTable[0][1]
		c.Info.Kd = c.GainTable[0][2]
	} else { // gain scheduling 2
		c.Info.Kp = c.GainTable[1][0]
		c.Info.Ki = c.GainTable[1][1]
		c.Info.Kd = c.GainTable[1][2]
	}

	// Proportional
	proportional := c.Info.Kp * err

	// Integrator (David page 49)
	c.Info.Integrator += DeltaTime * err
	integrator := c.Info.Integrator * c.Info.Ki

	// Differentiator (David page 49)
	differentiator := c.Info.Kd * (err - c.Info.PreviousError) / DeltaTime

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

func (c *GainPIDController) UpdateErrorSquare(p []float64) float64 {

	r := p[0] // goal
	y := p[1] // arrival rate

	// errors
	err := r - y

	// Proportional
	proportional := c.Info.Kp * err

	// Integrator (David page 49)
	c.Info.Integrator += DeltaTime * err
	integrator := c.Info.Integrator * c.Info.Ki

	// Differentiator (David page 49)
	differentiator := c.Info.Kd * (err - c.Info.PreviousError) / DeltaTime

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

func (c *GainPIDController) UpdateIncrementalForm(p []float64) float64 {

	r := p[0] // goal
	y := p[1] // arrival rate

	// errors
	err := r - y

	// Proportional // page 106
	//proportional := c.Info.Kp * err

	// Integrator // page 106
	c.Info.Integrator += DeltaTime * err
	integrator := c.Info.Integrator * c.Info.Ki

	// Differentiator // page 106
	//differentiator := c.Info.Kd * (err - c.Info.PreviousError) / DELTA_TIME

	// Delta of the new PC
	deltaU := c.Info.Kp*(err-c.Info.PreviousError) + c.Info.Ki*err*DeltaTime + c.Info.Kd*(err-2*c.Info.PreviousError+c.Info.PreviousPreviousError)/DeltaTime

	// pid output  TODO
	c.Info.Out = deltaU + integrator

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

func (c *GainPIDController) UpdateSmoothingDerivative(p []float64) float64 {

	r := p[0] // goal
	y := p[1] // arrival rate

	// errors
	err := r - y

	// Proportional
	proportional := c.Info.Kp * err

	// Integrator (David page 49)
	c.Info.Integrator += DeltaTime * err
	integrator := c.Info.Integrator * c.Info.Ki

	// smoothing the derivative term (page 104)
	differentiator := ALPHA*(err-c.Info.PreviousError)/DeltaTime + (1-ALPHA)*c.Info.PreviousDifferentiator
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

func (c *GainPIDController) SetGains(kp, ki, kd float64) {
	c.Info.Kp = kp
	c.Info.Ki = ki
	c.Info.Kd = kd
}

func (c *GainPIDController) SetKP(kp float64) {
	c.Info.Kp = kp
}

func (c *GainPIDController) GetKP() float64 {
	return c.Info.Kp
}
