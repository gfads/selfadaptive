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

func (Analyser) Run(fromMonitor chan shared.Symptoms, toPlanner chan shared.ToPlannerChan, goal string) {
	for {
		// receive symptom from monitor
		symptoms := <-fromMonitor
		info := shared.ToPlannerChan{}

		switch goal {
		case shared.NoWorry:
			info.ChangeRequest = shared.NoChange
		case shared.AlwaysSecure:
			switch symptoms.SecuritySymptom {
			case shared.Secure:
				info.ChangeRequest = shared.UseWeakCryptography
			case shared.Suspicious:
				info.ChangeRequest = shared.UseMediumCryptography
			case shared.Unsecure:
				info.ChangeRequest = shared.UseStrongCryptography
			default:
				shared.ErrorHandler(shared.GetFunction(), "Unknown environment symptom '"+symptoms.SecuritySymptom+"' unknown")
			}
		}
		toPlanner <- info
	}
}
