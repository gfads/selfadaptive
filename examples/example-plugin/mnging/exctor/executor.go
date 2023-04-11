/*********************************************************************************
Author: Nelson S Rosa
Description: This program implements a simple executor of MAPE-K.
Date: 28/02/2023
*********************************************************************************/

package exctor

import (
	"selfadaptive/examples/example-plugin/mnging/knwldge"
	"selfadaptive/shared"
)

type Executor struct{}

func NewExecutor() *Executor {
	return &Executor{}
}

func (Executor) Run(fromPlanner chan shared.ToPlannerChan, toManaged chan shared.ToManagedChan) {
	for {
		info := shared.ToManagedChan{}
		request := <-fromPlanner
		switch request.ChangeRequest {
		case shared.NoChange:
			info.Behaviours = knwldge.KnowledgeDatabase.AvailableBehaviours
			info.SelectedBehaviour = knwldge.KnowledgeDatabase.LastBehaviour
		case shared.UsePlainText:
			info.Behaviours = knwldge.KnowledgeDatabase.AvailableBehaviours
			info.SelectedBehaviour = "DefaultBehaviour"
		case shared.UseWeakCryptography:
			info.Behaviours = knwldge.KnowledgeDatabase.AvailableBehaviours
			info.SelectedBehaviour = "WeakCryptography"
		case shared.UseMediumCryptography:
			info.Behaviours = knwldge.KnowledgeDatabase.AvailableBehaviours
			info.SelectedBehaviour = "MediumCryptography"
		case shared.UseStrongCryptography:
			info.Behaviours = knwldge.KnowledgeDatabase.AvailableBehaviours
			info.SelectedBehaviour = "StrongCryptography"
		}

		knwldge.KnowledgeDatabase.LastBehaviour = info.SelectedBehaviour
		knwldge.KnowledgeDatabase.CurrentSecurityLevelOfApplication = info.SelectedBehaviour
		toManaged <- info
	}
}
