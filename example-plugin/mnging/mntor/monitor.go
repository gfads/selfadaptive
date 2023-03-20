package monitor

import (
	"selfadaptive/shared"
)

type Monitor struct{}

func NewMonitor() *Monitor {
	return &Monitor{}
}

func (Monitor) Start(fromManaged chan []func(), toAnalyser chan []func()) {
	for {
		// receive plugin behaviours
		pluginBehaviours := shared.LoadFuncs(shared.SourcesDir, shared.ExecutablesDir)

		// receive hard coded behaviours
		hardBehaviours := <-fromManaged

		// put all behaviours together
		var allBehaviours = []func(){}
		for i := 0; i < len(hardBehaviours); i++ {
			allBehaviours = append(allBehaviours, hardBehaviours[i])
		}

		for i := 0; i < len(pluginBehaviours); i++ {
			allBehaviours = append(allBehaviours, pluginBehaviours[i])
		}

		// send all behaviours to analyser
		toAnalyser <- allBehaviours
	}
}
