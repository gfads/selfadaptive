package main

import (
	"adaptive/adaptiveV5/environment/plugins/injector"
	"adaptive/adaptiveV5/selfadaptivesystem/managing"
	"adaptive/adaptiveV5/sharedadaptive"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	// Configure probes of the mntor
	var probes []func(managing.MonitorImpl)
	probes = append(probes, managing.MonitorImpl.ProbeSourceRemote)
	m := managing.NewMonitor(probes) // no mnged systems

	// Configure MAPE-K
	mapek := managing.NewMAPEK(m, nil, nil, nil)              // only mntor
	managingSystem := managing.NewManagingSystem(nil, &mapek) // no mnged system

	// Empty & initialise repositories
	inj := injector.PluginInjector{}
	inj.Initialize()

	// Start elements
	wg.Add(2)
	go inj.Start(sharedadaptive.REMOTE, &wg)
	go managingSystem.Start(&wg)

	//go AdaptationGoals()

	wg.Wait()
}
