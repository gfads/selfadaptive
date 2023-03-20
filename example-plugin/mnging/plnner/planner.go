package plnner

import "selfadaptive/shared"

type Planner struct{}

func NewPlanner() *Planner {
	return &Planner{}
}

func (Planner) Start(fromAnalyser chan shared.TypeChanManaging, toExecutor chan shared.TypeChanManaging) {
	for {
		info := <-fromAnalyser
		toExecutor <- info
	}
}
