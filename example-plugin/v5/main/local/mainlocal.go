package main

import (
	"adaptive/adaptiveV5/environment/plugins/injector"
	"adaptive/adaptiveV5/selfadaptivesystem/managed"
	"adaptive/adaptiveV5/selfadaptivesystem/managing"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	// Configure MAPE-K (no mntor)
	mapek := managing.NewMAPEK(nil, managing.NewAnalyser(), managing.NewPlanner(), managing.NewExecutor())
	managedSystem := managed.NewManaged()

	// Configure mnging system
	managingSystem := managing.NewManagingSystem(managedSystem, &mapek)

	// Empty repositories
	inj := injector.PluginInjector{}
	inj.Initialize()

	// Start mnged and mnging systems
	wg.Add(2)
	go managedSystem.Start(&wg)
	go managingSystem.Start(&wg)

	//go AdaptationGoals()   // TODO

	wg.Wait()
}
