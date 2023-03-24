/*********************************************************************************
Author: Nelson S Rosa
Description: This program implements a simple executor of MAPE-K.
Date: 28/02/2023
*********************************************************************************/

package exctor

import (
	"selfadaptive/example-plugin/mnging/knwldge"
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
		case shared.UsePlainText:
			info.Behaviours = knwldge.KnowledgeDatabase.AvailableBehaviours
			info.SelectedBehaviour = 0
		case shared.UseStrongCryptography:
			info.Behaviours = knwldge.KnowledgeDatabase.AvailableBehaviours
			info.SelectedBehaviour = 3
			/*case shared.KeepSecurity:
				info.Behaviours = knwldge.KnowledgeDatabase.AvailableBehaviours
				info.SelectedBehaviour = knwldge.KnowledgeDatabase.CurrentSecurityLevelOfApplication
			case shared.ImproveSecurity:
				info.Behaviours = knwldge.KnowledgeDatabase.AvailableBehaviours
				info.SelectedBehaviour = knwldge.KnowledgeDatabase.CurrentSecurityLevelOfApplication + 1 // TODO
			case shared.ReduceSecurity:
				info.Behaviours = knwldge.KnowledgeDatabase.AvailableBehaviours
				info.SelectedBehaviour = knwldge.KnowledgeDatabase.CurrentSecurityLevelOfApplication // TODO
			*/

		}

		knwldge.KnowledgeDatabase.LastBehaviour = info.SelectedBehaviour
		knwldge.KnowledgeDatabase.CurrentSecurityLevelOfApplication = info.SelectedBehaviour
		toManaged <- info
	}
}
