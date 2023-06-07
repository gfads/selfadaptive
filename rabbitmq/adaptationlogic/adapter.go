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
const SizeOfSameLevel = 60

var IncreasingGoal = []float64{500, 1000.0, 1250.0, 1500.0, 1750.0, 2000.0, 2250.0, 2500.0, 2750.0, 3000.0}
var RandomGoal = []float64{500.0, 900.0, 1250.0, 1300.0, 300.0, 800.0, 1000.0, 400.0, 500.0, 1000.0, 2000.0, 1500.0, 500.0, 1100.0, 1500.0, 600.0, 3000.0}
var L = 1.0
var tau = 1.0
var T = 0.1

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
	PC       int
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
	i := TrainingInfo{Kp: info.Kp, Ki: info.Ki, Kd: info.Kd, Data: []AdjustmenstInfo{}, TypeName: info.TypeName, SetPoint: setPoint, PC: pc}

	return AdaptationLogic{ExecutionType: executionType, TrainingInfo: i, FromBusiness: chFromBusiness, ToBusiness: chToBusiness, Controller: c, SetPoint: setPoint, MonitorInterval: monitorInterval * time.Second, PC: pc}
}

func (al AdaptationLogic) Run() {

	switch al.ExecutionType {
	case shared.StaticGoal:
		al.RunStaticGoal()
	case shared.DynamicGoal:
		al.RunDynamicGoal()
	case shared.RootLocusTraining:
		al.RootLocusTraining()
	case shared.ZieglerTraining:
		al.ZieglerTraining()
	case shared.CohenTraining:
		//al.ZieglerTraining() // TODO
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

func (al AdaptationLogic) RootLocusTrainingNewNotWorkingProperly() {
	fromAdjuster := make(chan TrainingInfo)
	toAdjuster := make(chan TrainingInfo)
	info := TrainingInfo{TypeName: al.TrainingInfo.TypeName}

	// create & execute adjustment mechanism
	go AdjustmentMechanism(fromAdjuster, toAdjuster, al.SetPoint)

	// warm up phase
	time.Sleep(WarmupTime * time.Second)

	// discard first measurement
	<-al.FromBusiness // receive no. of messages from business
	al.ToBusiness <- al.PC

	// loop of adaptation logic
	idx := 0
	for {

		n := <-al.FromBusiness // receive no. of messages from business

		// calculate new arrival rate (msg/s)
		rate := float64(n) / al.MonitorInterval.Seconds()
		info.Data = append(info.Data, AdjustmenstInfo{PC: al.PC, Rate: rate})

		if idx > 0 {
			info = CalculateRootLocusGains(info)
			if info.Kp > 0 && info.Ki > 0 { // end of training
				fmt.Printf("\"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)
				fmt.Println("***** End of Training - Copy/paste Gains *****")
				time.Sleep(10 * time.Hour)
				break
			} else {
				al.TrainingInfo.Kp = info.Kp
				al.TrainingInfo.Ki = info.Ki
				al.TrainingInfo.Kd = info.Kd
				fmt.Printf("\"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)
			}
		}
		// next pc training
		idx++
		al.PC += 1

		// send pc to business
		al.ToBusiness <- al.PC
	}
}

func (al AdaptationLogic) RootLocusTraining() {
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

		if l >= 1 && rate < (1.1*info.Data[l-1].Rate) { // TODO remove 1.10
			if l > TrainingSampleSize || tAttempts >= TrainingAttempts { // training is over
				info = CalculateRootLocusGains(info)
				al.TrainingInfo.Kp = info.Kp
				al.TrainingInfo.Ki = info.Ki
				al.TrainingInfo.Kd = info.Kd

				fmt.Printf("\"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)

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

func (al AdaptationLogic) ZieglerTraining() {
	//fromAdjuster := make(chan TrainingInfo)
	//toAdjuster := make(chan TrainingInfo)
	count := 0
	info := TrainingInfo{TypeName: al.TrainingInfo.TypeName}
	PCS := []int{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2}

	// create & execute adjustment mechanism
	//go AdjustmentMechanism(fromAdjuster, toAdjuster, al.SetPoint)

	// warm up phase
	time.Sleep(WarmupTime * time.Second)

	// discard first measurement
	<-al.FromBusiness       // receive no. of messages from business
	al.ToBusiness <- PCS[0] // Configure PC to execute the Ziegler Steps

	// loop of adaptation logic
	for {

		n := <-al.FromBusiness // receive no. of messages from business

		// calculate new arrival rate (msg/s)
		rate := float64(n) / al.MonitorInterval.Seconds()

		fmt.Println(al.PC, ";", rate)
		i := AdjustmenstInfo{PC: al.PC, Rate: rate}
		info.Data = append(info.Data, i)

		if count < len(PCS)-1 {
			count++
			// configure next PC
			al.PC = PCS[count]

		} else { // training is over
			info = CalculateZieglerGains(info)
			al.TrainingInfo.Kp = info.Kp
			al.TrainingInfo.Ki = info.Ki
			al.TrainingInfo.Kd = info.Kd

			fmt.Printf("Ziegler: \"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)

			info = CalculateCohenGains(info)
			al.TrainingInfo.Kp = info.Kp
			al.TrainingInfo.Ki = info.Ki
			al.TrainingInfo.Kd = info.Kd

			fmt.Printf("Cohen: \"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)

			info = CalculateAMIGOGains(info)
			al.TrainingInfo.Kp = info.Kp
			al.TrainingInfo.Ki = info.Ki
			al.TrainingInfo.Kd = info.Kd

			fmt.Printf("AMIGO: \"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)

			if al.TrainingInfo.Kp > 0 && al.TrainingInfo.Ki > 0 {
				fmt.Println("Bad gains...")
			}
			fmt.Println("***** End of Training - Copy/paste Gains *****")
			time.Sleep(10 * time.Hour)

		}

		// send pc to business
		al.ToBusiness <- al.PC
	}
}

func (al AdaptationLogic) CohenTraining() {
	//fromAdjuster := make(chan TrainingInfo)
	//toAdjuster := make(chan TrainingInfo)
	count := 0
	info := TrainingInfo{TypeName: al.TrainingInfo.TypeName}
	PCS := []int{5, 6, 5, 6, 5, 6, 5, 6, 5, 6, 5, 6}

	// create & execute adjustment mechanism
	//go AdjustmentMechanism(fromAdjuster, toAdjuster, al.SetPoint)

	// warm up phase
	time.Sleep(WarmupTime * time.Second)

	// discard first measurement
	<-al.FromBusiness       // receive no. of messages from business
	al.ToBusiness <- PCS[0] // Configure PC to execute the Ziegler Steps

	// loop of adaptation logic
	for {

		n := <-al.FromBusiness // receive no. of messages from business

		// calculate new arrival rate (msg/s)
		rate := float64(n) / al.MonitorInterval.Seconds()

		fmt.Println(al.PC, ";", rate)
		i := AdjustmenstInfo{PC: al.PC, Rate: rate}
		info.Data = append(info.Data, i)

		if count < len(PCS)-1 {
			count++
			// configure next PC
			al.PC = PCS[count]

		} else { // training is over
			info = CalculateCohenGains(info)
			al.TrainingInfo.Kp = info.Kp
			al.TrainingInfo.Ki = info.Ki
			al.TrainingInfo.Kd = info.Kd

			fmt.Printf("Cohen: \"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)

			if al.TrainingInfo.Kp > 0 && al.TrainingInfo.Ki > 0 {
				fmt.Println("Bad gains...")
			}
			fmt.Println("***** End of Training - Copy/paste Gains *****")
			time.Sleep(10 * time.Hour)

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
					info = CalculateRootLocusGains(info)
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
				info = CalculateRootLocusGains(info) // TODO
				toAdapter <- info                    // send new gains to adapter
			}
		}
	}
}

func CalculateRootLocusGains(info TrainingInfo) TrainingInfo {

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
		s1 += math.Pow(yLine[i], 2.0)
		s2 += uLine[i] * yLine[i]
		s3 += math.Pow(uLine[i], 2.0)
		s4 += yLine[i] * yLine[i+1]
		s5 += uLine[i] * yLine[i+1]
	}

	a := (s3*s4 - s2*s5) / (s1*s3 - math.Pow(s2, 2.0))
	b := (s1*s5 - s2*s4) / (s1*s3 - math.Pow(s2, 2.0))

	kp := 0.0
	ki := 0.0
	kd := 0.0

	//fmt.Printf("a=%.8f b=%.8f\n", a, b)

	switch info.TypeName {
	case shared.BasicP:
		kp = (1 + a) / b // Feedback Control of Computing Systems - page 264
		ki = 0.0
		kd = 0.0
	case shared.BasicPi:
		kp = (a - 0.36) / b
		ki = (a - b*kp) / b
		kd = 0.0
	case shared.BasicPid:
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

func CalculateZieglerGains(info TrainingInfo) TrainingInfo {

	sumRate1 := 0.0
	sumRate2 := 0.0

	for i := 2; i < len(info.Data); i++ { // discard 2 initial measurements, i.e., 5 samples
		if i%2 == 0 {
			sumRate1 += info.Data[i].Rate
		} else {
			sumRate2 += info.Data[i].Rate
		}
	}

	fmt.Println("************")
	dataSize := float64((len(info.Data) - 2) / 2.0) // discard 2 initial measurements
	meanRate1 := sumRate1 / dataSize
	meanRate2 := sumRate2 / dataSize

	diffRate := meanRate2 - meanRate1
	K := diffRate

	switch info.TypeName {
	case shared.BasicP:
		info.Kp = tau / K * L
		info.Ki = 0.0
		info.Kd = 0.0
	case shared.BasicPi:
		ti := L / 0.3
		info.Kp = 0.9 * tau / K * L
		info.Ki = info.Kp / ti
		info.Kd = 0.0
	case shared.BasicPid:
		ti := 2 * L
		td := 0.5 * L
		info.Kp = 1.2 * tau / K * L
		info.Ki = info.Kp / ti
		info.Kd = info.Kp * td
	}
	return info
}

func CalculateCohenGains(info TrainingInfo) TrainingInfo {

	sumRate1 := 0.0
	sumRate2 := 0.0

	for i := 2; i < len(info.Data); i++ { // discard 2 initial measurements, i.e., 5 samples
		if i%2 == 0 {
			sumRate1 += info.Data[i].Rate
		} else {
			sumRate2 += info.Data[i].Rate
		}
	}

	fmt.Println("************")
	dataSize := float64((len(info.Data) - 2) / 2.0) // discard 2 initial measurements
	meanRate1 := sumRate1 / dataSize
	meanRate2 := sumRate2 / dataSize

	diffRate := meanRate2 - meanRate1

	K := diffRate
	theta := tau / (tau + T)

	switch info.TypeName {
	case shared.BasicP:
		info.Kp = 1 / K * (1 + (0.35 * theta / (1 - theta))) * T / tau
		info.Ki = 0.0
		info.Kd = 0.0
	case shared.BasicPi:
		ti := ((3.3 - 3.0*theta) / (1 + 1.2*theta)) * tau
		info.Kp = 0.9 / K * (1 + (0.92 * theta / (1 - theta))) * T / tau
		info.Ki = info.Kp / ti
		info.Kd = 0.0
	case shared.BasicPid:
		ti := ((2.5 - 2.0*theta) / (1 - 0.39*theta)) * tau
		td := ((0.37 * (1 - theta)) / (1 - 0.81*theta)) * tau
		info.Kp = 1.35 / K * (1 + (0.18 * theta / (1 - theta))) * T / tau
		info.Ki = info.Kp / ti
		info.Kd = info.Kp * td
	}

	return info
}

func CalculateAMIGOGains(info TrainingInfo) TrainingInfo {

	sumRate1 := 0.0
	sumRate2 := 0.0

	for i := 2; i < len(info.Data); i++ { // discard 2 initial measurements, i.e., 5 samples
		if i%2 == 0 {
			sumRate1 += info.Data[i].Rate
		} else {
			sumRate2 += info.Data[i].Rate
		}
	}

	fmt.Println("************")
	dataSize := float64((len(info.Data) - 2) / 2.0) // discard 2 initial measurements
	meanRate1 := sumRate1 / dataSize
	meanRate2 := sumRate2 / dataSize

	diffRate := meanRate2 - meanRate1
	K := diffRate

	switch info.TypeName {
	case shared.BasicPi:
		ti := (0.35 + 13*math.Pow(T, 2)/(math.Pow(T, 2)+12*tau*T+7*math.Pow(tau, 2))) * tau
		info.Kp = 1 / K * (0.15 + (0.35-tau*T/math.Pow(tau+T, 2))*T/tau)
		info.Ki = info.Kp / ti
		info.Kd = 0.0
	case shared.BasicPid:
		ti := ((0.4*tau + 0.8*T) / (tau + 0.1*T)) * tau
		td := (0.5 * T / (0.3*tau + T)) * tau
		info.Kp = 1 / K * (0.2 + 0.45*T/tau)
		info.Ki = info.Kp / ti
		info.Kd = info.Kp * td
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
