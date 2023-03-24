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
	"selfadaptive/shared"
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
		case shared.SecureEnvironment:
			info.ChangeRequest = shared.UsePlainText
		case shared.UnsecureEnvironment:
			info.ChangeRequest = shared.UseStrongCryptography
		default:
			shared.ErrorHandler(shared.GetFunction(), "Unknown environment symptom '"+symptoms.SecuritySymptom+"' unknown")
		}

		// configure and send request change to planner
		toPlanner <- info
	}
}
