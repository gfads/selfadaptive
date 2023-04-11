package mnging

import (
	"selfadaptive/examples/example-plugin/envrnment"
	"selfadaptive/examples/example-plugin/mnging/anlser"
	"selfadaptive/examples/example-plugin/mnging/exctor"
	"selfadaptive/examples/example-plugin/mnging/mntor"
	"selfadaptive/examples/example-plugin/mnging/plnner"
	"selfadaptive/shared"
)

type ManagingSystem struct {
	Goal string
}

func NewManagingSystem(g string) *ManagingSystem {
	return &ManagingSystem{Goal: g}
}

func (ms ManagingSystem) Run(fromManaged chan map[string]func(), toManaged chan shared.ToManagedChan, env *envrnment.Environment) {

	toAnalyser := make(chan shared.Symptoms)
	toPlanner := make(chan shared.ToPlannerChan)
	toExecutor := make(chan shared.ToPlannerChan)

	m := mntor.NewMonitor()
	a := anlser.NewAnalyser()
	p := plnner.NewPlanner()
	e := exctor.NewExecutor()

	go m.Run(fromManaged, toAnalyser, env)
	go a.Run(toAnalyser, toPlanner, ms.Goal)
	go p.Run(toPlanner, toExecutor)
	go e.Run(toExecutor, toManaged)
}
