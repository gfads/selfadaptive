package anlser

import (
	"math/rand"
	"selfadaptive/shared"
)

type Analyser struct{}

func NewAnalyser() *Analyser {
	return &Analyser{}
}

func (Analyser) Start(fromMonitor chan []func(), toManaged chan shared.TypeChanManaging) {

	info := shared.TypeChanManaging{}
	for {

		// receive behaviours from mntor
		allBehaviours := <-fromMonitor

		// configure ans send info to managed system
		info.Functions = allBehaviours
		info.N = rand.Intn(len(allBehaviours))

		toManaged <- info
	}
}
