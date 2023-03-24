package mntor

import (
	"selfadaptive/example-plugin/envrnment"
	"selfadaptive/example-plugin/mnging/knwldge"
	"selfadaptive/shared"
	"strconv"
	"time"
)

type Monitor struct{}

func NewMonitor() *Monitor {
	return &Monitor{}
}

func (Monitor) Run(fromManaged chan map[string]func(), toAnalyser chan shared.Symptoms) {
	for {
		// monitor interval
		time.Sleep(5 * time.Second)

		// sense environment
		securityLevel, pluginBehaviours := envrnment.NewEnvironment().Sense()

		// sense managed system
		hardBehaviours := <-fromManaged

		// put all behaviours together
		var allBehaviours = map[string]func(){}
		for i := range hardBehaviours {
			allBehaviours[i] = hardBehaviours[i]
		}

		for i := 0; i < len(pluginBehaviours); i++ {
			allBehaviours["Plugin"+strconv.Itoa(i)] = pluginBehaviours[i]
		}

		// generate symptom
		symptoms := shared.Symptoms{}

		// update available behaviours symptom
		if len(allBehaviours) > len(knwldge.KnowledgeDatabase.AvailableBehaviours) {
			symptoms.PluginSymptom = shared.NewPluginvAvailable
		} else {
			symptoms.PluginSymptom = shared.NoNewPluginAvailable
		}

		// update security symptom
		symptoms.SecuritySymptom = securityLevel

		// update knowledge database
		knwldge.KnowledgeDatabase.AvailableBehaviours = allBehaviours

		// send all behaviours to analyser
		toAnalyser <- symptoms
	}
}
