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
	"selfadaptive/example-plugin/mnging/knwldge"
	"selfadaptive/shared"
	"strconv"
)

type Analyser struct{}

func NewAnalyser() *Analyser {
	return &Analyser{}
}

func (Analyser) Run(fromMonitor chan shared.Symptoms, toPlanner chan shared.ToPlannerChan, goal int) {
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
			info.ChangeRequest = shared.ImproveSecurity
		default:
			shared.ErrorHandler(shared.GetFunction(), "Unknown environment symptom '"+strconv.Itoa(symptoms.SecuritySymptom)+"' unknown")
		}

		// configure and send request change to planner
		toPlanner <- info
	}
}
