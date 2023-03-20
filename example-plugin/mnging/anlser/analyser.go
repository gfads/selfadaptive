/*********************************************************************************
Author: Nelson S Rosa
Description: This code implements an analyser that has several strategies for
selecting the following behaviour of the managed system:
V0: No adaptation
V1: randomly select the subsequent behaviour.
V2: always use the behaviour of the last plugin made available.
Date: 14/02/2023
*********************************************************************************/
package anlser

import (
	"math/rand"
	"selfadaptive/shared"
)

type Analyser struct{}

func NewAnalyser() *Analyser {
	return &Analyser{}
}

func (Analyser) Start(fromMonitor chan []func(), toManaged chan shared.TypeChanManaging, goal string) {
	switch goal {
	case shared.NoAdaptation:
		startV0(fromMonitor, toManaged)
	case shared.AnyBehaviour:
		startV1(fromMonitor, toManaged)
	case shared.AlwaysUpdated:
		startV2(fromMonitor, toManaged)
	default:
		shared.ErrorHandler(shared.GetFunction(), "Analyser unknown")
	}
}

// No Adaptation
func startV0(fromMonitor chan []func(), toManaged chan shared.TypeChanManaging) {
	info := shared.TypeChanManaging{}
	for {

		// receive behaviours from monitor
		allBehaviours := <-fromMonitor

		// configure and send info to managed system
		info.Functions = allBehaviours
		info.N = 0

		toManaged <- info
	}
}

// Any behaviour
func startV1(fromMonitor chan []func(), toManaged chan shared.TypeChanManaging) {
	info := shared.TypeChanManaging{}
	for {

		// receive behaviours from mntor
		allBehaviours := <-fromMonitor

		// configure and send info to managed system
		info.Functions = allBehaviours
		info.N = rand.Intn(len(allBehaviours))

		toManaged <- info
	}
}

// Always updated
func startV2(fromMonitor chan []func(), toManaged chan shared.TypeChanManaging) {

	info := shared.TypeChanManaging{}

	for {

		// receive behaviours from monitor
		allBehaviours := <-fromMonitor

		// configure and send info to managed system
		info.Functions = allBehaviours
		info.N = len(allBehaviours) - 1

		// send info to managed system
		toManaged <- info
	}
}
