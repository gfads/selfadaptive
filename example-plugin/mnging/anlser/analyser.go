/*********************************************************************************
Author: Nelson S Rosa
Description: This code implements an analyser that has several strategies for
selecting the following behaviour of the managed system:
No adaptation: use same behaviour as before
Any: randomly select the subsequent behaviour.
Always updated: always use the behaviour of the last plugin made available.
Date: 19/03/2023
*********************************************************************************/
package anlser

import (
	"math/rand"
	"selfadaptive/example-plugin/mnging/knwldge"
	"selfadaptive/shared"
	"strconv"
)

type Analyser struct{}

func NewAnalyser() *Analyser {
	return &Analyser{}
}

func (Analyser) Run(fromMonitor chan shared.Symptoms, toPlanner chan shared.ToPlannerChan, goal string) {
	for {

		// receive symptom from monitor
		symptoms := <-fromMonitor
		info := shared.ToPlannerChan{}

		// analyse of symptoms (* priority to security *)
		switch symptoms.SecuritySymptom {
		case shared.HighSecureEnvironment:
			info.ChangeRequest = shared.UsePlainText
		case shared.MediumSecureEnvironment:
			if knwldge.KnowledgeDatabase.CurrentSecurityLevelOfEnvironment > shared.MediumSecureEnvironment {
				info.ChangeRequest = shared.ReduceSecurity
			}
			if knwldge.KnowledgeDatabase.CurrentSecurityLevelOfEnvironment == shared.MediumSecureEnvironment {
				info.ChangeRequest = shared.KeepSecurity
			}
			if knwldge.KnowledgeDatabase.CurrentSecurityLevelOfEnvironment < shared.MediumSecureEnvironment {
				info.ChangeRequest = shared.ImproveSecurity
			}
		case shared.LowSecureEnvironment:
			if knwldge.KnowledgeDatabase.CurrentSecurityLevelOfEnvironment > shared.MediumSecureEnvironment {
				info.ChangeRequest = shared.ReduceSecurity
			}
			if knwldge.KnowledgeDatabase.CurrentSecurityLevelOfEnvironment == shared.MediumSecureEnvironment {
				info.ChangeRequest = shared.KeepSecurity
			}
			if knwldge.KnowledgeDatabase.CurrentSecurityLevelOfEnvironment < shared.MediumSecureEnvironment {
				info.ChangeRequest = shared.ImproveSecurity
			}
		}
		switch symptoms.PluginSymptom {
		case shared.NoNewBehaviourAvailable:
			switch goal {
			case shared.NoAdaptation:
				info.ChangeRequest = shared.NoChangeRequest
				info.SelectedBehaviour = knwldge.KnowledgeDatabase.LastBehaviour
			case shared.AnyBehaviour:
				info.ChangeRequest = shared.AnyBehaviourRequest
				info.SelectedBehaviour = rand.Intn(len(knwldge.KnowledgeDatabase.AvailableBehaviours))
			case shared.AlwaysUpdated:
				info.ChangeRequest = shared.LastBehaviourRequest
				info.SelectedBehaviour = knwldge.KnowledgeDatabase.LastBehaviour
			default:
				shared.ErrorHandler(shared.GetFunction(), "Goal '"+goal+"' unknown")
			}
		case shared.NewBehaviourAvailable:
			switch goal {
			case shared.NoAdaptation:
				info.ChangeRequest = shared.NoChangeRequest
				info.SelectedBehaviour = knwldge.KnowledgeDatabase.LastBehaviour
			case shared.AnyBehaviour:
				info.ChangeRequest = shared.AnyBehaviourRequest
				info.SelectedBehaviour = rand.Intn(len(knwldge.KnowledgeDatabase.AvailableBehaviours))
			case shared.AlwaysUpdated:
				info.ChangeRequest = shared.LastBehaviourRequest
				info.SelectedBehaviour = len(knwldge.KnowledgeDatabase.AvailableBehaviours) - 1
			default:
				shared.ErrorHandler(shared.GetFunction(), "Goal '"+goal+"' unknown")
			}
		default:
			shared.ErrorHandler(shared.GetFunction(), "Symptom '"+strconv.Itoa(symptom)+"'unknown")
		}

		// configure and send request change to planner
		toPlanner <- info
	}
}
