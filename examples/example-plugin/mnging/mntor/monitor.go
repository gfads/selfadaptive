package mntor

import (
	"main.go/examples/example-plugin/envrnment"
	"main.go/examples/example-plugin/mnging/knwldge"
	"main.go/shared"
	"strconv"
	"time"
)

type Monitor struct{}

func NewMonitor() *Monitor {
	return &Monitor{}
}

func (Monitor) Run(fromManaged chan map[string]func(), toAnalyser chan shared.Symptoms, env *envrnment.Environment) {
	for {
		// monitor interval
		time.Sleep(5 * time.Second)

		// sense environment
		securityLevel, pluginBehaviours := env.Sense()

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

		// define symptom
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
		knwldge.KnowledgeDatabase.CurrentSecurityLevelOfEnvironment = securityLevel

		// send all behaviours to analyser
		toAnalyser <- symptoms
	}
}
