package adaptationlogic

import (
	"fmt"
	"main.go/controllers/def/info"
	"main.go/controllers/def/ops"
	"main.go/shared"
	"math"
	"os"
	"time"
)

const TrainingSampleSize = 60
const TimeBetweenAdjustments = 1200 // seconds
const MaximumNrmse = 0.30
const WarmupTime = 30 // seconds
const TrainingAttempts = 30
const SizeOfSameLevel = 40

var IncreasingGoal = []float64{500, 1000.0, 1250.0, 1500.0, 1750.0, 2000.0, 2250.0, 2500.0, 2750.0, 3000.0}
var RandomGoal = []float64{500, 200.0, 1250.0, 1500.0, 800, 2000.0, 2250.0, 350.0, 480.0, 3000.0}

type AdjustmenstInfo struct {
	PC   int
	Rate float64
}

type TrainingInfo struct {
	TypeName string
	Kp       float64
	Ki       float64
	Kd       float64
	Data     []AdjustmenstInfo
	SetPoint float64
}

type AdaptationLogic struct {
	MonitorInterval time.Duration
	FromBusiness    chan int
	ToBusiness      chan int
	Controller      ops.IController
	SetPoint        float64
	PC              int
	ExecutionType   string
	TrainingInfo    TrainingInfo
}

func NewAdaptationLogic(executionType string, chFromBusiness chan int, chToBusiness chan int, info info.Controller, setPoint float64, monitorInterval time.Duration, pc int) AdaptationLogic {

	c := ops.NewController(info)
	i := TrainingInfo{Kp: info.Kp, Ki: info.Ki, Kd: info.Kd, Data: []AdjustmenstInfo{}, TypeName: info.TypeName, SetPoint: setPoint}

	return AdaptationLogic{ExecutionType: executionType, TrainingInfo: i, FromBusiness: chFromBusiness, ToBusiness: chToBusiness, Controller: c, SetPoint: setPoint, MonitorInterval: monitorInterval * time.Second, PC: pc}
}

func (al AdaptationLogic) Run() {

	switch al.ExecutionType {
	case shared.StaticGoal:
		al.RunStaticGoal()
	case shared.DynamicGoal:
		al.RunDynamicGoal()
	case shared.OffLineTraining:
		al.RunOfflineTraining()
	case shared.OnLineTraining:
		al.RunOnlineTraining()
	default:
		fmt.Println("Execution type ´", al.ExecutionType, "´ is invalid")
		os.Exit(0)
	}
}

func (al AdaptationLogic) RunDynamicGoal() {

	// discard first measurement
	<-al.FromBusiness // receive no. of messages from business
	al.ToBusiness <- al.PC

	currentGoalIdx := 0
	count := 0

	// loop of adaptation logic
	for {
		select {
		case n := <-al.FromBusiness: // interaction with the business
			count++

			// calculate new arrival rate (msg/s)
			rate := float64(n) / al.MonitorInterval.Seconds()

			// catch pc and its yielded rate
			fmt.Println(al.PC, ";", rate, ";", RandomGoal[currentGoalIdx])

			// invoke controller to calculate new pc
			pc := int(math.Round(al.Controller.Update(RandomGoal[currentGoalIdx], rate)))

			// update pc at adaptation mechanism
			al.PC = pc

			// send new pc to business
			al.ToBusiness <- al.PC

			// change goal
			if count >= SizeOfSameLevel {
				count = 0
				currentGoalIdx++
				if currentGoalIdx >= len(RandomGoal) {
					fmt.Println("********** Copy/paste data of the experiments **********")
					time.Sleep(10 * time.Hour)
				}
			}
		}
	}
}

func (al AdaptationLogic) RunStaticGoal() {

	// discard first measurement
	<-al.FromBusiness // receive no. of messages from business
	al.ToBusiness <- al.PC

	// loop of adaptation logic
	for {
		select {
		case n := <-al.FromBusiness: // interaction with the business

			// calculate new arrival rate (msg/s)
			rate := float64(n) / al.MonitorInterval.Seconds()

			// catch pc and its yielded rate
			fmt.Println(al.PC, ";", rate, ";", al.SetPoint)

			// invoke controller to calculate new pc
			pc := int(math.Round(al.Controller.Update(al.SetPoint, rate)))

			// update pc at adaptation mechanism
			al.PC = pc

			// send new pc to business
			al.ToBusiness <- al.PC
		}
	}
}

func (al AdaptationLogic) RunOfflineTraining() {
	fromAdjuster := make(chan TrainingInfo)
	toAdjuster := make(chan TrainingInfo)
	tAttempts := 0
	info := TrainingInfo{TypeName: al.TrainingInfo.TypeName}

	// create & execute adjustment mechanism
	go AdjustmentMechanism(fromAdjuster, toAdjuster, al.SetPoint)

	// warm up phase
	time.Sleep(WarmupTime * time.Second)

	// discard first measurement
	<-al.FromBusiness // receive no. of messages from business
	al.ToBusiness <- al.PC

	// loop of adaptation logic
	for {

		n := <-al.FromBusiness // receive no. of messages from business

		// calculate new arrival rate (msg/s)
		rate := float64(n) / al.MonitorInterval.Seconds()

		l := len(info.Data)

		if l >= 1 && rate < info.Data[l-1].Rate {
			if l > TrainingSampleSize || tAttempts >= TrainingAttempts { // training is over
				info = CalculateGains(info)
				al.TrainingInfo.Kp = info.Kp
				al.TrainingInfo.Ki = info.Ki
				al.TrainingInfo.Kd = info.Kd

				fmt.Printf("Kp= %.8f Ki=%.8f  Kd=%.8f\n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)

				if al.TrainingInfo.Kp > 0 && al.TrainingInfo.Ki > 0 {
					fmt.Println("Bad gains...")
				}
				fmt.Println("***** End of Training - Copy/paste Gains *****")
				time.Sleep(10 * time.Hour)

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
	}
}

func (al AdaptationLogic) RunOnlineTraining() {
	state := 0
	fromAdjuster := make(chan TrainingInfo)
	toAdjuster := make(chan TrainingInfo)
	tAttempts := 0
	countCycles := 0
	countExperiments := 0

	info := TrainingInfo{TypeName: al.TrainingInfo.TypeName}

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
					al.TrainingInfo.Kp = info.Kp
					al.TrainingInfo.Ki = info.Ki
					al.TrainingInfo.Kd = info.Kd

					// initial configuration of controller
					al.Controller.SetGains(al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)

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
				pc := int(math.Round(al.Controller.Update(al.SetPoint, rate)))

				// update pc at adaptation mechanism
				al.PC = pc

				// send new pc to business
				al.ToBusiness <- al.PC

				// update set point dynamically -- ONLY FOR EXPERIMENTS
				if countCycles > SizeOfSameLevel {
					al.SetPoint = RandomGoal[countExperiments]
					countCycles = 0
					countExperiments++
				} else {
					countCycles++
				}

			case toAdjuster <- info: // interaction with the adjustment mechanism
				info = <-fromAdjuster

				al.TrainingInfo.Kp = info.Kp
				al.TrainingInfo.Ki = info.Ki
				al.TrainingInfo.Kd = info.Kd

				al.Controller.SetGains(al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)

				// reset data
				info.Data = []AdjustmenstInfo{}
			default:
			}
		}
	}
}

func (al *AdaptationLogic) UpdateSetPoint(n float64) {
	al.SetPoint = n
}

func AdjustmentMechanism(toAdapter chan TrainingInfo, fromAdapter chan TrainingInfo, setPoint float64) {
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

func CalculateGains(info TrainingInfo) TrainingInfo {

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

	if info.TypeName == "P" { // TODO
		kp = (1 + a) / b
	}

	if info.TypeName == "PI" { // TODO
		kp = (a - 0.36) / b
		ki = (a - b*kp) / b
	}

	if info.TypeName == shared.BasicPid {
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

	return info
}

func CalculateNRMSE(info TrainingInfo) float64 {

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
