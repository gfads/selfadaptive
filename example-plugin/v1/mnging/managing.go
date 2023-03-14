package mnging

import (
	"selfadaptive/example-plugin/v1/mnging/anlser"
	"selfadaptive/example-plugin/v1/mnging/mntor"
	"selfadaptive/shared/channeltypes"
)

type ManagingSystem struct{}

func NewManagingSystem() *ManagingSystem {
	return &ManagingSystem{}
}

func (m ManagingSystem) Start(fromManaged chan []func(), toManaged chan channeltypes.TypeChanManaging) {

	//fromMonitor := make(chan []func())
	toAnalyser := make(chan []func())

	analyser := anlser.NewAnalyser()
	monitor := evolutive.NewMonitor()

	go analyser.Start(toAnalyser, toManaged)
	go monitor.Start(fromManaged, toAnalyser)
}
