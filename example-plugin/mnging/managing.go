package mnging

import (
	"selfadaptive/example-plugin/mnging/anlser"
	"selfadaptive/example-plugin/mnging/mntor"
	"selfadaptive/shared"
)

type ManagingSystem struct {
	Goal string
}

func NewManagingSystem(g string) *ManagingSystem {
	return &ManagingSystem{Goal: g}
}

func (m ManagingSystem) Start(fromManaged chan []func(), toManaged chan shared.TypeChanManaging) {

	//fromMonitor := make(chan []func())
	toAnalyser := make(chan []func())

	analyser := anlser.NewAnalyser()
	monitor := monitor.NewMonitor()

	go analyser.Start(toAnalyser, toManaged, m.Goal)
	go monitor.Start(fromManaged, toAnalyser)
}
