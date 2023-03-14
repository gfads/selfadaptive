package main

import (
	"adaptive/adaptiveV4/environment/plugins/injector"
	"adaptive/adaptiveV4/selfadaptivesystem/managed"
	"adaptive/adaptiveV4/selfadaptivesystem/managing"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	// Configure probes of the mntor
	//var probes []func(mnging.MonitorImpl) mnging.MonitorInfo
	//probes = append(probes, mnging.MonitorImpl.ProbeSourceRemote)
	//m := mnging.NewMonitor(probes) // no mnged systems

	// Configure MAPE-K (no mntor)
	mapek := managing.NewMAPEK(nil, managing.NewAnalyser(), managing.NewPlanner(), managing.NewExecutor())
	managedSystem := managed.NewManaged()

	// Configure mnging system
	managingSystem := managing.NewManagingSystem(managedSystem, &mapek)

	// Empty repositories
	inj := injector.PluginInjector{}
	inj.Initialize()

	// Start elements
	wg.Add(2)
	go managedSystem.Start(&wg)
	go managingSystem.Start(&wg)

	//go AdaptationGoals()

	wg.Wait()
}
