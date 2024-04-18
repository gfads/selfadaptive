package mnging

import (
	"main.go/examples/example-plugin/envrnment"
	"main.go/examples/example-plugin/mnging/anlser"
	"main.go/examples/example-plugin/mnging/exctor"
	"main.go/examples/example-plugin/mnging/mntor"
	"main.go/examples/example-plugin/mnging/plnner"
	"main.go/shared"
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
