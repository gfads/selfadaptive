package adaptationlogic

import (
	"fmt"
	"math"
	"selfadaptive/controllers/def/info"
	"selfadaptive/controllers/def/ops"
	"time"
)

type AdaptationLogic struct {
	MonitorInterval time.Duration
	FromBusiness    chan int
	ToBusiness      chan int
	Controller      ops.IController
	SetPoint        float64
	PC              int
}

func NewAdaptationLogic(chFromBusiness chan int, chToBusiness chan int, info info.Controller, setPoint float64, monitorInterval time.Duration, pc int) AdaptationLogic {
	c := ops.NewController(info)
	return AdaptationLogic{FromBusiness: chFromBusiness, ToBusiness: chToBusiness, Controller: c, SetPoint: setPoint, MonitorInterval: monitorInterval * time.Second, PC: pc}
}

func (al AdaptationLogic) Run() {

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
