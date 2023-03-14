package evolutive

import (
	"selfadaptive/shared"
)

type Monitor struct{}

const SourcesDir = "/Volumes/GoogleDrive/Meu Drive/go/selfadaptive/example-plugin/v1/env/plugins/source"
const ExecutablesDir = "/Volumes/GoogleDrive/Meu Drive/go/selfadaptive/example-plugin/v1/env/plugins/executable"

func NewMonitor() *Monitor {
	return &Monitor{}
}

func (Monitor) Start(fromManaged chan []func(), toAnalyser chan []func()) {
	for {
		// receive plugin behaviours
		pluginBehaviours := shared.LoadFuncs(SourcesDir, ExecutablesDir)

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
