package adaptationlogic

import (
	"fmt"
	"github.com/montanaflynn/stats"
	"main.go/controllers/def/ops"
	"main.go/controllers/def/parameters"
	"main.go/shared"
	"math"
	"os"
	"time"
)

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
	FromBusiness    chan shared.SubscriberToAdapter
	ToBusiness      chan int
	Controller      ops.IController
	SetPoint        float64
	PC              int
	ExecutionType   string
	TrainingInfo    TrainingInfo
	File            *os.File
}

func NewAdaptationLogic(chFromBusiness chan shared.SubscriberToAdapter, chToBusiness chan int, p parameters.ExecutionParameters, df *os.File) AdaptationLogic {

	c := ops.NewController(p)
	i := TrainingInfo{
		Kp:       *p.Kp,
		Ki:       *p.Ki,
		Kd:       *p.Kd,
		Data:     []AdjustmenstInfo{},
		TypeName: *p.ControllerType,
		SetPoint: *p.SetPoint,
		PC:       *p.PrefetchCount}

	return AdaptationLogic{ExecutionType: *p.ExecutionType,
		TrainingInfo:    i,
		FromBusiness:    chFromBusiness,
		ToBusiness:      chToBusiness,
		Controller:      c,
		SetPoint:        *p.SetPoint,
		MonitorInterval: time.Duration(*p.MonitorInterval) * time.Second,
		PC:              *p.PrefetchCount,
		File:            df}
}

func (al AdaptationLogic) Run() {

	switch al.ExecutionType {
	case shared.ExperimentalDesign:
		al.RunExperimentalDesign()
	case shared.Experiment:
		al.RunExperiment()
	case shared.RootTraining:
		al.RootLocusTraining()
	case shared.ZieglerTraining: // all non-analytical tune methods
		al.ZieglerTraining()
	//case shared.CohenTraining:
	//al.ZieglerTraining() // TODO
	case shared.OnLineTraining:
		al.RunOnlineTraining()
	case shared.WebTraining:
		al.PIDTunerWeb()
	default:
		fmt.Println("Execution type ´", al.ExecutionType, "´ is invalid")
		os.Exit(0)
	}
}

func (al AdaptationLogic) RunExperimentalDesign() {
	count := 0
	info := TrainingInfo{TypeName: al.TrainingInfo.TypeName}

	// warm up phase
	time.Sleep(shared.WarmupTime * time.Second)

	// discard first measurement
	<-al.FromBusiness // receive no. of messages from business
	al.ToBusiness <- al.PC

	// loop of adaptation logic
	idx := 0
	for {
		n := <-al.FromBusiness // receive no. of messages from business

		// calculate new arrival rate (msg/s)
		rate := float64(n.ReceivedMessages) / n.D

		if count > shared.SizeOfSameLevel {
			count = 0
			// calculate mean
			meanPc, meanRate := calculateMeans(info)

			fmt.Printf("%v;%v;%.6f;%.6f\n", n.QueueSize, meanPc, meanRate, 0.0)
			fmt.Fprintf(al.File, "%v;%v;%.6f;%.6f\n", 0.0, meanPc, meanRate, 0.0)

			// reset info
			info = TrainingInfo{TypeName: al.TrainingInfo.TypeName}

			// update pc
			idx++
			al.PC = shared.InputSteps[idx]

			// check end of experiment
			if idx > len(shared.InputSteps) {
				os.Exit(0) // end of experiment
			}
		} else {
			count++

			//fmt.Printf("%v;%.6f\n", al.PC, rate)
			// store training pc and rate
			info.Data = append(info.Data, AdjustmenstInfo{PC: al.PC, Rate: rate})
		}

		// send pc to business
		al.ToBusiness <- al.PC
	}
}

func (al AdaptationLogic) RunExperiment() {

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
			rate := float64(n.ReceivedMessages) / n.D

			// catch pc and its yielded rate
			fmt.Fprintf(al.File, "%d;%d;%f;%f\n", n.QueueSize, al.PC, rate, shared.RandomGoals[currentGoalIdx])
			fmt.Printf("%d;%d;%f;%f\n", n.QueueSize, al.PC, rate, shared.RandomGoals[currentGoalIdx])

			// invoke controller to calculate new pc
			if al.TrainingInfo.TypeName == shared.Fuzzy {

				// calculate pc update
				update := int(math.Round(al.Controller.Update(shared.RandomGoals[currentGoalIdx], rate)))

				// update pc at adaptation mechanism
				al.PC = al.PC + update
				if al.PC <= 0 {
					al.PC = 1
				}
			} else {
				pc := int(math.Round(al.Controller.Update(shared.RandomGoals[currentGoalIdx], rate)))
				al.PC = pc
			}

			// send new pc to business
			al.ToBusiness <- al.PC

			// check the need for changing the setpoint
			if count >= shared.SizeOfSameLevel {
				count = 0
				currentGoalIdx++
				if currentGoalIdx >= len(shared.RandomGoals) {
					al.File.Close()
					fmt.Println("Adapter", al.File.Name())
					os.Exit(0)
				}
			}
		}
	}
}

// Just generate data to be used at pidtuner.com - no calculation of control gains is made
func (al AdaptationLogic) PIDTunerWeb() {

	var totalTime float64

	// warm up phase
	time.Sleep(shared.WarmupTime * time.Second)

	// discard first measurement
	<-al.FromBusiness // receive no. of messages from business
	al.ToBusiness <- al.PC

	// loop of adaptation logic
	nSameLevel := 1
	for {

		n := <-al.FromBusiness // receive no. of messages from business

		// calculate new arrival rate (msg/s)
		rate := float64(n.ReceivedMessages) / al.MonitorInterval.Seconds()

		fmt.Println(totalTime, ";", al.PC, ";", rate)

		// next pc training
		nSameLevel++
		totalTime += al.MonitorInterval.Seconds()

		if nSameLevel > 30 {
			nSameLevel = 0
			al.PC += 1
		}

		if totalTime > 5000 {
			fmt.Println("End of experiments - copy and paste data")
			time.Sleep(10 * time.Hour)
		}

		// send pc to business
		al.ToBusiness <- al.PC
	}
}

func (al AdaptationLogic) RootLocusTrainingNewNotWorkingProperly() {
	fromAdjuster := make(chan TrainingInfo)
	toAdjuster := make(chan TrainingInfo)
	info := TrainingInfo{TypeName: al.TrainingInfo.TypeName}

	// create & execute adjustment mechanism
	go AdjustmentMechanism(fromAdjuster, toAdjuster, al.SetPoint)

	// warm up phase
	time.Sleep(shared.WarmupTime * time.Second)

	// discard first measurement
	<-al.FromBusiness // receive no. of messages from business
	al.ToBusiness <- al.PC

	// loop of adaptation logic
	idx := 0
	for {

		n := <-al.FromBusiness // receive no. of messages from business

		// calculate new arrival rate (msg/s)
		rate := float64(n.ReceivedMessages) / al.MonitorInterval.Seconds()
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
	time.Sleep(shared.WarmupTime * time.Second)

	// discard first measurement
	<-al.FromBusiness // receive no. of messages from business
	al.ToBusiness <- al.PC

	rates := []float64{}
	pcs := []float64{}
	correlation := 1.0
	err := error(nil)

	// loop of adaptation logic
	for {

		n := <-al.FromBusiness // receive no. of messages from business

		// calculate new arrival rate (msg/s)
		rate := float64(n.ReceivedMessages) / al.MonitorInterval.Seconds()

		// store measured data
		rates = append(rates, rate)
		pcs = append(pcs, float64(al.PC))

		l := len(info.Data)
		if l > 2 { // mimimum of 3 values
			correlation, err = stats.Correlation(pcs, rates)
			if err != nil {
				shared.ErrorHandler("Correlation error", shared.GetFunction())
				return
			}
		}
		//if l >= 1 && rate < (1.1*info.Data[l-1].Rate) { // TODO remove 1.10
		//if l >= 1 && rate < ((1+1.0/float64(al.PC+1))*info.Data[l-1].Rate) { // TODO remove
		//fmt.Println("Correlation index= ", correlation, "Sampling size= ", len(rates), rates)
		if l > 2 && correlation < 0.995 {
			if l > shared.TrainingSampleSize || tAttempts >= shared.TrainingAttempts { // training is over
				info = CalculateRootLocusGains(info)
				al.TrainingInfo.Kp = info.Kp
				al.TrainingInfo.Ki = info.Ki
				al.TrainingInfo.Kd = info.Kd

				fmt.Printf("\n \"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)
				fmt.Fprintf(al.File, "\n -kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)
				if al.TrainingInfo.Kp > 0 && al.TrainingInfo.Ki > 0 {
					fmt.Println("Bad gains...")
				}
				al.File.Close()
				fmt.Println("***** End of Root Locus Tuning *****")
				os.Exit(0)

			} else { // redo the training && keep already collected data
				//fmt.Println("Training attempts=", tAttempts, len(info.Data), rate, al.PC)
				tAttempts++

				// remove last data used to calculate correlation
				pcs = pcs[:len(pcs)-1]
				rates = rates[:len(rates)-1]
			}
		} else { // continue the training
			tAttempts = 0
			fmt.Printf("%v;%v;%f;%f\n", n.QueueSize, al.PC, rate, 0.0)
			fmt.Fprintf(al.File, "%v;%v;%f;%f\n", n.QueueSize, al.PC, rate, 0.0)

			// store training pc and rate
			info.Data = append(info.Data, AdjustmenstInfo{PC: al.PC, Rate: rate})

			// increment pc
			al.PC += 1
		}

		// send pc to business
		al.ToBusiness <- al.PC
	}
}

func (al AdaptationLogic) ZieglerTrainingOld() {
	count := 0
	info := TrainingInfo{TypeName: al.TrainingInfo.TypeName}

	// warm up phase
	//time.Sleep(shared.WarmupTime * time.Second)

	// discard first measurement
	<-al.FromBusiness                     // receive no. of messages from business
	al.ToBusiness <- shared.InputSteps[0] // Configure PC to execute the Ziegler Steps

	// loop of adaptation logic
	for {
		count++
		// receive no. of messages from business
		n := <-al.FromBusiness

		// calculate new arrival rate (msg/s)
		rate := float64(n.ReceivedMessages) / al.MonitorInterval.Seconds()

		fmt.Printf("%v;%v;%f;%f\n", n.QueueSize, al.PC, rate, 0.0)
		fmt.Fprintf(al.File, "%v;%v;%f;%f\n", n.QueueSize, al.PC, rate, 0.0)

		d := AdjustmenstInfo{PC: al.PC, Rate: rate}
		info.Data = append(info.Data, d)

		if count <= len(shared.InputSteps) {
			if count < len(shared.InputSteps)/2 { // TODO
				al.PC = shared.InputSteps[0]
			} else {
				al.PC = shared.InputSteps[1]
			}
		} else { // training is over
			info = CalculateZieglerGains(info)
			al.TrainingInfo.Kp = info.Kp
			al.TrainingInfo.Ki = info.Ki
			al.TrainingInfo.Kd = info.Kd

			fmt.Printf("Ziegler: \"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)
			fmt.Fprintf(al.File, "Ziegler: \"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)

			info = CalculateCohenGains(info)
			al.TrainingInfo.Kp = info.Kp
			al.TrainingInfo.Ki = info.Ki
			al.TrainingInfo.Kd = info.Kd

			fmt.Printf("Cohen: \"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)
			fmt.Fprintf(al.File, "Cohen: \"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)

			info = CalculateAMIGOGains(info)
			al.TrainingInfo.Kp = info.Kp
			al.TrainingInfo.Ki = info.Ki
			al.TrainingInfo.Kd = info.Kd

			fmt.Printf("AMIGO: \"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)
			fmt.Fprintf(al.File, "AMIGO: \"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)

			if al.TrainingInfo.Kp > 0 && al.TrainingInfo.Ki > 0 {
				fmt.Println("Bad gains...")
			}
			fmt.Println("***** End of Training - Copy/paste Gains *****")
			al.File.Close()
			os.Exit(0)
		}

		// send pc to business
		al.ToBusiness <- al.PC
	}
}

func (al AdaptationLogic) ZieglerTraining() {
	info := TrainingInfo{TypeName: al.TrainingInfo.TypeName}

	// warm up phase
	//time.Sleep(shared.WarmupTime * time.Second)

	// discard first measurement
	<-al.FromBusiness                     // receive no. of messages from business
	al.ToBusiness <- shared.InputSteps[0] // Configure PC to execute the Ziegler Steps

	// loop of adaptation logic
	for i := 0; i < shared.ZieglerRepetitions; i++ {
		for j := 0; j < shared.SizeOfSameLevelZiegler; j++ {
			// receive no. of messages from business
			n := <-al.FromBusiness

			// calculate new arrival rate (msg/s)
			rate := float64(n.ReceivedMessages) / al.MonitorInterval.Seconds()

			fmt.Printf("%v;%v;%f;%f\n", n.QueueSize, al.PC, rate, 0.0)
			fmt.Fprintf(al.File, "%v;%v;%f;%f\n", n.QueueSize, al.PC, rate, 0.0)

			d := AdjustmenstInfo{PC: al.PC, Rate: rate}
			info.Data = append(info.Data, d)

			// send pc to business
			al.ToBusiness <- al.PC
		}
		if i%2 == 0 {
			al.PC = shared.InputStepsZiegler[0]
		} else {
			al.PC = shared.InputStepsZiegler[1]
		}
	}

	// training is over
	info = CalculateZieglerGains(info)
	al.TrainingInfo.Kp = info.Kp
	al.TrainingInfo.Ki = info.Ki
	al.TrainingInfo.Kd = info.Kd

	fmt.Printf("Ziegler: \"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)
	fmt.Fprintf(al.File, "Ziegler: \"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)

	info = CalculateCohenGains(info)
	al.TrainingInfo.Kp = info.Kp
	al.TrainingInfo.Ki = info.Ki
	al.TrainingInfo.Kd = info.Kd

	fmt.Printf("Cohen: \"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)
	fmt.Fprintf(al.File, "Cohen: \"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)

	info = CalculateAMIGOGains(info)
	al.TrainingInfo.Kp = info.Kp
	al.TrainingInfo.Ki = info.Ki
	al.TrainingInfo.Kd = info.Kd

	fmt.Printf("AMIGO: \"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)
	fmt.Fprintf(al.File, "AMIGO: \"-kp=%.8f\", \"-ki=%.8f\", \"-kd=%.8f\" \n", al.TrainingInfo.Kp, al.TrainingInfo.Ki, al.TrainingInfo.Kd)

	al.File.Close()
	os.Exit(0)
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
	//time.Sleep(shared.WarmupTime * time.Second)

	// discard first measurement
	<-al.FromBusiness       // receive no. of messages from business
	al.ToBusiness <- PCS[0] // Configure PC to execute the Ziegler Steps

	// loop of adaptation logic
	for {

		n := <-al.FromBusiness // receive no. of messages from business

		// calculate new arrival rate (msg/s)
		rate := float64(n.ReceivedMessages) / al.MonitorInterval.Seconds()

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
	//time.Sleep(shared.WarmupTime * time.Second)

	// discard first measurement
	<-al.FromBusiness // receive no. of messages from business
	al.ToBusiness <- al.PC

	// loop of adaptation logic
	for {
		switch state {

		case 0: // training phase
			n := <-al.FromBusiness // receive no. of messages from business

			// calculate new arrival rate (msg/s)
			rate := float64(n.ReceivedMessages) / al.MonitorInterval.Seconds()

			l := len(info.Data)

			if l >= 1 && rate < info.Data[l-1].Rate {
				if l > shared.TrainingSampleSize || tAttempts >= shared.TrainingAttempts { // training is over
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
				rate := float64(n.ReceivedMessages) / al.MonitorInterval.Seconds()

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
				if countCycles > shared.SizeOfSameLevel {
					al.SetPoint = shared.RandomGoals[countExperiments]
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
		time.Sleep(shared.TimeBetweenAdjustments * time.Second) // wait for xx seconds before next adjusting
		info := <-fromAdapter
		nrmse := CalculateNRMSE(info)
		switch state {
		case 0:
			if nrmse < shared.MaximumNrmse { /// TODO
				fmt.Printf("No update of control gains %.4f < %.4f %d\n", nrmse, shared.MaximumNrmse, len(info.Data))
				toAdapter <- info // send new gains
				break             // previous gain improved rate - nothing to do
			} else { // recalculate gains
				fmt.Printf("Update control gains %.4f >= %.4f %d\n", nrmse, shared.MaximumNrmse, len(info.Data))
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

	fmt.Printf("a=%.8f b=%.8f\n", a, b)

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

	sumDeltaY1 := 0.0
	sumDeltaY2 := 0.0

	for i := 0; i < len(info.Data); i++ { // discard 2 initial input steps, i.e., 5 samples
		if i < len(shared.InputSteps)/2 {
			sumDeltaY1 += info.Data[i].Rate
		} else {
			sumDeltaY2 += info.Data[i].Rate
		}
	}

	fmt.Println("************")
	dataSize := float64((len(info.Data) - 2) / 2.0) // discard 2 initial input steps
	k1 := sumDeltaY1 / dataSize
	k2 := sumDeltaY2 / dataSize
	K := k2 - k1
	lambda := K * shared.Tau / shared.T // see Figure 9-2

	switch info.TypeName {
	case shared.BasicP:
		k := 1 / lambda
		info.Kp = k
		info.Ki = 0.0
		info.Kd = 0.0
	case shared.BasicPi:
		k := 0.9 / lambda
		Ti := 3.0 * shared.Tau
		info.Kp = k
		info.Ki = k / Ti
		info.Kd = 0.0
	case shared.BasicPid:
		k := 1.2 / lambda
		Ti := 2.0 * shared.Tau
		Td := shared.Tau / 2.0
		info.Kp = k
		info.Ki = k / Ti
		info.Kd = k * Td
	}
	return info
}

func CalculateCohenGains(info TrainingInfo) TrainingInfo {

	sumDeltaY1 := 0.0
	sumDeltaY2 := 0.0

	for i := 0; i < len(info.Data); i++ { // discard 2 initial input steps, i.e., 5 samples
		if i < len(shared.InputSteps)/2 {
			sumDeltaY1 += info.Data[i].Rate
		} else {
			sumDeltaY2 += info.Data[i].Rate
		}
	}

	fmt.Println("************")
	dataSize := float64((len(info.Data) - 2) / 2.0) // discard 2 initial measurements
	k1 := sumDeltaY1 / dataSize
	k2 := sumDeltaY2 / dataSize
	K := k2 - k1
	theta := shared.Tau / (shared.Tau + shared.T)

	switch info.TypeName {
	case shared.BasicP:
		k := 1 / K * (1 + (0.35 * theta / (1 - theta))) * shared.T / shared.Tau
		info.Kp = k
		info.Ki = 0.0
		info.Kd = 0.0
	case shared.BasicPi:
		k := 0.9 / K * (1 + (0.92 * theta / (1 - theta))) * shared.T / shared.Tau
		Ti := ((3.3 - 3.0*theta) / (1 + 1.2*theta)) * shared.Tau
		info.Kp = k
		info.Ki = k / Ti
		info.Kd = 0.0
	case shared.BasicPid:
		k := 1.35 / K * (1 + (0.18 * theta / (1 - theta))) * shared.T / shared.Tau
		Ti := ((2.5 - 2.0*theta) / (1 - 0.39*theta)) * shared.Tau
		Td := (0.37 * (1 - theta) / (1 - 0.81*theta)) * shared.Tau
		info.Kp = k
		info.Ki = k / Ti
		info.Kd = k * Td
	}

	return info
}

func CalculateAMIGOGains(info TrainingInfo) TrainingInfo {

	sumDeltaY1 := 0.0
	sumDeltaY2 := 0.0

	for i := 0; i < len(info.Data); i++ { // discard 2 initial input steps, i.e., 5 samples
		if i < len(shared.InputSteps)/2 {
			sumDeltaY1 += info.Data[i].Rate
		} else {
			sumDeltaY2 += info.Data[i].Rate
		}
	}

	fmt.Println("************")
	dataSize := float64((len(info.Data) - 2) / 2.0) // discard 2 initial measurements
	k1 := sumDeltaY1 / dataSize
	k2 := sumDeltaY2 / dataSize
	K := k2 - k1

	switch info.TypeName {
	case shared.BasicPi:
		k := 1 / K * (0.15 + (0.35-shared.Tau*shared.T/math.Pow(shared.Tau+shared.T, 2))*shared.T/shared.Tau)
		Ti := (0.35 + 13*math.Pow(shared.T, 2)/(math.Pow(shared.T, 2)+12*shared.Tau*shared.T+7*math.Pow(shared.Tau, 2))) * shared.Tau
		info.Kp = k
		info.Ki = k / Ti
		info.Kd = 0.0
	case shared.BasicPid:
		k := 1 / K * (0.2 + 0.45*shared.T/shared.Tau)
		Ti := ((0.4*shared.Tau + 0.8*shared.T) / (shared.Tau + 0.1*shared.T)) * shared.Tau
		Td := (0.5 * shared.T / (0.3*shared.Tau + shared.T)) * shared.Tau
		info.Kp = k
		info.Ki = k / Ti
		info.Kd = k * Td
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

func calculateMeans(d TrainingInfo) (int, float64) {
	sPc := 0
	sRate := 0.0
	for i := 0; i < len(d.Data); i++ {
		sPc += d.Data[i].PC
		sRate += d.Data[i].Rate
	}
	meanPC := sPc / len(d.Data)
	meanRate := sRate / float64(len(d.Data))

	return meanPC, meanRate
}
