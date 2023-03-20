package exctor

import "selfadaptive/shared"

type Executor struct{}

func NewExecutor() *Executor {
	return &Executor{}
}

func (Executor) Start(fromPlanner chan shared.TypeChanManaging, toManaged chan shared.TypeChanManaging) {
	for {
		info := <-fromPlanner
		toManaged <- info
	}
}
