package mnging

import (
	"selfadaptive/example-plugin/mnging/anlser"
	"selfadaptive/example-plugin/mnging/exctor"
	"selfadaptive/example-plugin/mnging/mntor"
	"selfadaptive/example-plugin/mnging/plnner"
	"selfadaptive/shared"
)

type ManagingSystem struct {
	Goal int
}

func NewManagingSystem(g int) *ManagingSystem {
	return &ManagingSystem{Goal: g}
}

func (m ManagingSystem) Run(fromManaged chan []func(), toManaged chan shared.ToManagedChan) {

	toAnalyser := make(chan shared.Symptoms)
	toPlanner := make(chan shared.ToPlannerChan)
	toExecutor := make(chan shared.ToPlannerChan)

	monitor := mntor.NewMonitor()
	analyser := anlser.NewAnalyser()
	planner := plnner.NewPlanner()
	executor := exctor.NewExecutor()

	go monitor.Run(fromManaged, toAnalyser)
	go analyser.Run(toAnalyser, toPlanner, m.Goal)
	go planner.Run(toPlanner, toExecutor)
	go executor.Run(toExecutor, toManaged)
}
