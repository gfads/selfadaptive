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
		case shared.NoChangeRequest:
			info.Behaviours = knwldge.KnowledgeDatabase.AvailableBehaviours
			info.SelectedBehaviour = request.SelectedBehaviour
		case shared.LastBehaviourRequest:
			info.Behaviours = knwldge.KnowledgeDatabase.AvailableBehaviours
			info.SelectedBehaviour = request.SelectedBehaviour
		case shared.AnyBehaviourRequest:
			info.Behaviours = knwldge.KnowledgeDatabase.AvailableBehaviours
			info.SelectedBehaviour = request.SelectedBehaviour
		}

		knwldge.KnowledgeDatabase.LastBehaviour = info.SelectedBehaviour
		toManaged <- info
	}
}
