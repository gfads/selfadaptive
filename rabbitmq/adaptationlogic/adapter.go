package adaptationlogic

import (
	"fmt"
	"math"
	"os"
	"selfadaptive/rabbitmq/controller"
	"time"
)

const TrainingSampleSize = 30
const TimeBetweenAdjustments = 1200 // seconds
const MaximumNrmse = 0.30
const WarmupTime = 30 // seconds
const TrainingAttempts = 10
const SizeOfSameLevel = 20

var PcExperimentSet = []float64{1000.0, 1250.0, 1500.0, 1750.0, 2000.0, 2250.0, 2500.0, 2750.0, 3000.0}

type AdjustmenstInfo struct {
	PC   int
	Rate float64
}

type TODOINFO struct {
	Kp             float64
	Ki             float64
	Kd             float64
	Data           []AdjustmenstInfo
	ControllerType string
	SetPoint       float64
	PIDType        string
}

type AdaptationLogic struct {
	MonitorInterval time.Duration
	FromBusiness    chan int
	ToBusiness      chan int
	ControllerType  string
	PIDType         string
	Kp              float64
	Ki              float64
	Kd              float64
	SetPoint        float64
	PC              int
	TrainingType    string
}

func NewAdaptationLogic(chFromBusiness chan int, chToBusiness chan int, controllerType string, pidType string, trainingType string, kp float64, ki float64, kd float64, setPoint float64, monitorInterval time.Duration, pc int) AdaptationLogic {
	return AdaptationLogic{FromBusiness: chFromBusiness, ToBusiness: chToBusiness, ControllerType: controllerType, PIDType: pidType, TrainingType: trainingType, Kp: kp, Ki: ki, Kd: kd, SetPoint: setPoint, MonitorInterval: monitorInterval * time.Second, PC: pc}
}

func (al AdaptationLogic) Run() {

	switch al.TrainingType {
	case "Offline":
		al.RunOfflineTraining()
	case "Online":
		al.RunOnlineTraining()
	default:
		fmt.Println("Training type ´", al.TrainingType, "´ is invalid")
		os.Exit(0)
	}
}
func (al AdaptationLogic) RunOnlineTraining() {
	state := 0
	fromAdjuster := make(chan TODOINFO)
	toAdjuster := make(chan TODOINFO)
	tAttempts := 0
	countCycles := 0
	countExperiments := 0

	info := TODOINFO{Kp: al.Kp, Ki: al.Ki, Kd: al.Kd, Data: []AdjustmenstInfo{}, ControllerType: al.ControllerType, SetPoint: al.SetPoint}

	// create controller
	var c controller.IController
	c = controller.NewController(al.ControllerType, al.PIDType, al.Kp, al.Ki, al.Kd)

	// create & execute adjustment mechanism
	go AdjustmentMechanism(fromAdjuster, toAdjuster, al.SetPoint)

	// warm up phase
	time.Sleep(WarmupTime * time.Second)

	// discard first measurement
	<-al.FromBusiness // receive no. of messages from business
	al.ToBusiness <- al.PC

	// loop of adaptation logic
	for {
		switch state {

		case 0: // training phase
			n := <-al.FromBusiness // receive no. of messages from business

			// calculate new arrival rate (msg/s)
			rate := float64(n) / al.MonitorInterval.Seconds()

			l := len(info.Data)

			if l >= 1 && rate < info.Data[l-1].Rate {
				if l > TrainingSampleSize || tAttempts >= TrainingAttempts { // training is over
					info = CalculateGains(info)
					al.Kp = info.Kp
					al.Ki = info.Ki
					al.Kd = info.Kd

					// initial configuration of controller
					c.SetGains(al.Kp, al.Ki, al.Kd)

					// reset data
					info.Data = []AdjustmenstInfo{}

					// restore initial pc
					al.PC = 1

					// set new state
					state = 1
				} else { // redo the training && keep already collected data
					//fmt.Println("Training attempts=", tAttempts, len(info.Data), rate, al.PC)
					tAttempts++
				}
			} else { // continue the training
				fmt.Println(al.PC, ";", rate)

				// store training pc and rate
				info.Data = append(info.Data, AdjustmenstInfo{PC: al.PC, Rate: rate})

				// increment pc
				al.PC += 1
			}

			// send pc to business
			al.ToBusiness <- al.PC
			break

		case 1: // regular execution
			select {
			case n := <-al.FromBusiness: // interaction with the business

				// calculate new arrival rate (msg/s)
				rate := float64(n) / al.MonitorInterval.Seconds()

				// store pc and rate
				info.Data = append(info.Data, AdjustmenstInfo{PC: al.PC, Rate: rate})

				// catch pc and its yielded rate
				fmt.Println(al.PC, ";", rate, ";", al.SetPoint)

				// invoke controller to calculate new pc
				pc := int(math.Round(c.Update(al.SetPoint, rate)))

				// update pc at adaptation mechanism
				al.PC = pc

				// send new pc to business
				al.ToBusiness <- al.PC

				// update set point dynamically -- ONLY FOR EXPERIMENTS
				if countCycles > SizeOfSameLevel {
					//al.SetPoint = myRandon()
					al.SetPoint = PcExperimentSet[countExperiments]
					countCycles = 0
					countExperiments++
				} else {
					countCycles++
				}

			case toAdjuster <- info: // interaction with the adjustment mechanism
				info = <-fromAdjuster

				al.Kp = info.Kp
				al.Ki = info.Ki
				al.Kd = info.Kd

				c.SetGains(al.Kp, al.Ki, al.Kd)

				// reset data
				info.Data = []AdjustmenstInfo{}
			default:
			}
		}
	}
}

func (al AdaptationLogic) RunOfflineTraining() {
	fromAdjuster := make(chan TODOINFO)
	toAdjuster := make(chan TODOINFO)
	countCycles := 0
	countExperiments := 0

	info := TODOINFO{Kp: al.Kp, Ki: al.Ki, Kd: al.Kd, Data: []AdjustmenstInfo{}, ControllerType: al.ControllerType, SetPoint: al.SetPoint}

	// create controller
	var c controller.IController
	c = controller.NewController(al.ControllerType, al.PIDType, al.Kp, al.Ki, al.Kd)

	// create & execute adjustment mechanism
	go AdjustmentMechanism(fromAdjuster, toAdjuster, al.SetPoint)

	// warm up phase
	time.Sleep(WarmupTime * time.Second)

	// discard first measurement
	<-al.FromBusiness // receive no. of messages from business
	al.ToBusiness <- al.PC

	// loop of adaptation logic
	for {
		select {
		case n := <-al.FromBusiness: // interaction with the business

			// calculate new arrival rate (msg/s)
			rate := float64(n) / al.MonitorInterval.Seconds()

			// store pc and rate
			info.Data = append(info.Data, AdjustmenstInfo{PC: al.PC, Rate: rate})

			// catch pc and its yielded rate
			fmt.Println(al.PC, ";", rate, ";", al.SetPoint)

			// invoke controller to calculate new pc
			pc := int(math.Round(c.Update(al.SetPoint, rate)))

			// update pc at adaptation mechanism
			al.PC = pc

			// send new pc to business
			al.ToBusiness <- al.PC

			// update set point dynamically -- ONLY FOR EXPERIMENTS
			if countCycles > SizeOfSameLevel {
				// update set point dynamically -- ONLY FOR EXPERIMENTS
				if countCycles > SizeOfSameLevel {
					//al.SetPoint = myRandon()
					al.SetPoint = PcExperimentSet[countExperiments]
					countCycles = 0
					countExperiments++
				} else {
					countCycles++
				}
			} else {
				countCycles++
			}

		case toAdjuster <- info: // interaction with the adjustment mechanism
			info = <-fromAdjuster

			al.Kp = info.Kp
			al.Ki = info.Ki
			al.Kd = info.Kd

			c.SetGains(al.Kp, al.Ki, al.Kd)

			// reset data
			info.Data = []AdjustmenstInfo{}
		default:
		}
	}
}

func (al *AdaptationLogic) UpdateSetPoint(n float64) {
	al.SetPoint = n
}

func AdjustmentMechanism(toAdapter chan TODOINFO, fromAdapter chan TODOINFO, setPoint float64) {
	state := 0

	for {
		time.Sleep(TimeBetweenAdjustments * time.Second) // wait for xx seconds before next adjusting
		info := <-fromAdapter
		nrmse := CalculateNRMSE(info)
		switch state {
		case 0:
			if nrmse < MaximumNrmse { /// TODO
				fmt.Printf("No update of control gains %.4f < %.4f %d\n", nrmse, MaximumNrmse, len(info.Data))
				toAdapter <- info // send new gains
				break             // previous gain improved rate - nothing to do
			} else { // recalculate gains
				fmt.Printf("Update control gains %.4f >= %.4f %d\n", nrmse, MaximumNrmse, len(info.Data))
				info = CalculateGains(info)
				toAdapter <- info // send new gains to adapter
			}
		}
	}
}

func CalculateGains(info TODOINFO) TODOINFO {

	// calculate mean
	sumU := 0.0
	sumY := 0.0
	for i := 0; i < len(info.Data)-1; i++ {
		sumU += float64(info.Data[i].PC)
	}

	for i := 1; i < len(info.Data); i++ {
		sumY += info.Data[i].Rate
	}

	mu := sumU / float64(len(info.Data)-1)
	my := sumY / float64(len(info.Data)-1)

	s1 := 0.0
	s2 := 0.0
	s3 := 0.0
	s4 := 0.0
	s5 := 0.0

	uLine := make([]float64, len(info.Data))
	yLine := make([]float64, len(info.Data))

	for i := 0; i < len(info.Data); i++ {
		uLine[i] = float64(info.Data[i].PC) - mu
		yLine[i] = info.Data[i].Rate - my
	}

	for i := 0; i < len(info.Data)-1; i++ {
		s1 += yLine[i] * yLine[i]
		s2 += uLine[i] * yLine[i]
		s3 += uLine[i] * uLine[i]
		s4 += yLine[i] * yLine[i+1]
		s5 += uLine[i] * yLine[i+1]
	}

	a := (s3*s4 - s2*s5) / (s1*s3 - s2*s2)
	b := (s1*s5 - s2*s4) / (s1*s3 - s2*s2)

	kp := 0.0
	ki := 0.0
	kd := 0.0

	if info.ControllerType == "P" {
		kp = (1 + a) / b
	}

	if info.ControllerType == "PI" {
		kp = (a - 0.36) / b
		ki = (a - b*kp) / b
	}

	if info.ControllerType == "PID" {
		kd = 0.11 / b
		kp = (-0.063 + a - 2*b*kd) / b
		ki = (0.3 - b*kp - b*kd + a) / b
	}

	// reconfigure gains
	if !math.IsNaN(a) && !math.IsNaN(b) {
		info.Kp = kp
		info.Ki = ki
		info.Kd = kd
	}

	fmt.Printf("Kp= %.8f Ki=%.8f  Kd=%.8f\n", info.Kp, info.Ki, info.Kd)

	// TODO - remove in the future
	if kp > 0 && ki > 0 {
		fmt.Println("Bad gains... EXIT")
		os.Exit(0)
	}
	return info
}

func CalculateNRMSE(info TODOINFO) float64 {

	minRate := info.SetPoint
	maxRate := info.SetPoint
	s := 0.0

	for i := 0; i < len(info.Data); i++ {
		// define min/max rates
		if info.Data[i].Rate > maxRate {
			maxRate = info.Data[i].Rate
		}
		if info.Data[i].Rate < minRate {
			minRate = info.Data[i].Rate
		}
		s += math.Pow(info.SetPoint-info.Data[i].Rate, 2.0)
	}
	rmse := math.Sqrt(s / float64(len(info.Data)))
	nmrse := rmse / (maxRate - minRate)

	fmt.Println("NMRSE:: ", nmrse, rmse, maxRate, minRate, len(info.Data))
	return nmrse
}
