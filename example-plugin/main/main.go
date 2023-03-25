/*********************************************************************************
Author: Nelson S Rosa
Description: This program randomly selects a behaviour between 2 plugins and 3
hard-coded behaviours.
Date: 04/02/2023
*********************************************************************************/
package main

import (
	"fmt"
	"selfadaptive/example-plugin/envrnment"
	"selfadaptive/example-plugin/mnged"
	"selfadaptive/example-plugin/mnging"
	"selfadaptive/shared"
)

func main() {

	// remove old plugins
	shared.RemoveContents(shared.ExecutablesDir)
	shared.RemoveContents(shared.SourcesDir)

	// configure the adaptation goal
	goal := shared.AlwaysSecure
	//goal := shared.NoWorry

	// instantiate channels
	fromManaged := make(chan map[string]func())
	toManaged := make(chan shared.ToManagedChan)

	// instantiate elements
	managed := mnged.NewManagedElement()
	managing := mnging.NewManagingSystem(goal)
	environment := envrnment.NewEnvironment()

	fmt.Println("******* GOAL: ", goal, "*******")
	go environment.Run()
	go managed.Run(fromManaged, toManaged)
	go managing.Run(fromManaged, toManaged, environment)

	_, _ = fmt.Scanln()
}
