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

func (m ManagingSystem) Start(fromManaged chan []func(), toManaged chan shared.TypeChanManaging) {

	toAnalyser := make(chan []func())
	toPlanner := make(chan shared.TypeChanManaging)
	toExecutor := make(chan shared.TypeChanManaging)

	monitor := mntor.NewMonitor()
	analyser := anlser.NewAnalyser()
	planner := plnner.NewPlanner()
	executor := exctor.NewExecutor()

	go monitor.Start(fromManaged, toAnalyser)
	go analyser.Start(toAnalyser, toPlanner, m.Goal)
	go planner.Start(toPlanner, toExecutor)
	go executor.Start(toExecutor, toManaged)
}
