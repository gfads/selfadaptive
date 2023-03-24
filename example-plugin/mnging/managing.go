package mnging

import (
	"selfadaptive/example-plugin/mnging/anlser"
	"selfadaptive/example-plugin/mnging/exctor"
	"selfadaptive/example-plugin/mnging/mntor"
	"selfadaptive/example-plugin/mnging/plnner"
	"selfadaptive/shared"
)

type ManagingSystem struct {
	Goal string
}

func NewManagingSystem(g string) *ManagingSystem {
	return &ManagingSystem{Goal: g}
}

func (ms ManagingSystem) Run(fromManaged chan map[string]func(), toManaged chan shared.ToManagedChan) {

	toAnalyser := make(chan shared.Symptoms)
	toPlanner := make(chan shared.ToPlannerChan)
	toExecutor := make(chan shared.ToPlannerChan)

	m := mntor.NewMonitor()
	a := anlser.NewAnalyser()
	p := plnner.NewPlanner()
	e := exctor.NewExecutor()

	go m.Run(fromManaged, toAnalyser)
	go a.Run(toAnalyser, toPlanner, ms.Goal)
	go p.Run(toPlanner, toExecutor)
	go e.Run(toExecutor, toManaged)
}
