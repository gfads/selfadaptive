package evolutive

import (
	"selfadaptive/shared"
)

type Monitor struct{}

func NewMonitor() *Monitor {
	return &Monitor{}
}

func (Monitor) Start(fromManaged chan []func(), toAnalyser chan []func()) {
	dirSources := shared.BaseDirPlugins + "/" + shared.Version + "/" + shared.SourcesDir
	dirExecutables := shared.BaseDirPlugins + "/" + shared.Version + "/" + shared.ExecutablesDir

	for {
		// receive plugin behaviours
		pluginBehaviours := shared.LoadFuncs(dirSources, dirExecutables)

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
