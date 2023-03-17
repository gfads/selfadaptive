/*********************************************************************************
Author: Nelson S Rosa
Description: This program randomly selects a behaviour between 2 plugins and 3
hard-coded behaviours.
Date: 04/02/2023
*********************************************************************************/
package main

import (
	"fmt"
	"selfadaptive/example-plugin/v1/envrnment"
	"selfadaptive/example-plugin/v1/mnged"
	"selfadaptive/example-plugin/v1/mnging"
	"selfadaptive/shared"
)

func main() {

	// configure version of the example
	shared.Version = "v1"

	// instantiate channels
	fromManaged := make(chan shared.TypeChanManaging)
	fromManaging := make(chan []func())

	// instantiate elements
	managed := mnged.NewManagedElement()
	managing := mnging.NewManagingSystem()
	environment := envrnment.NewEnvironment()

	//
	go environment.Start()
	go managed.Start(fromManaging, fromManaged)
	go managing.Start(fromManaging, fromManaged)

	_, _ = fmt.Scanln()
}
