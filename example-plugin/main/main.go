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

	// configure the adaptation goal
	//goal := shared.AlwaysUpdated
	//goal := shared.AnyBehaviour
	goal := shared.NoAdaptation

	// instantiate channels
	fromManaged := make(chan shared.TypeChanManaging)
	fromManaging := make(chan []func())

	// instantiate elements
	managed := mnged.NewManagedElement()
	managing := mnging.NewManagingSystem(goal)
	environment := envrnment.NewEnvironment()

	//
	go environment.Start()
	go managed.Start(fromManaging, fromManaged)
	go managing.Start(fromManaging, fromManaged)

	_, _ = fmt.Scanln()
}
