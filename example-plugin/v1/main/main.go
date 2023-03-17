/*********************************************************************************
Author: Nelson S Rosa
Description: This program randomly selects a behaviour between 2 plugins and 3
hard-coded behaviours.
Date: 04/02/2023
*********************************************************************************/
package main

import (
	"fmt"
	"selfadaptive/example-plugin/v1/env"
	"selfadaptive/example-plugin/v1/mnged"
	"selfadaptive/example-plugin/v1/mnging"
	"selfadaptive/shared"
)

func main() {

	fromManaged := make(chan shared.TypeChanManaging)
	fromManaging := make(chan []func())

	managed := mnged.NewManagedElement()
	managing := mnging.NewManagingSystem()
	environment := env.NewEnvironment()

	go environment.Start()
	go managed.Start(fromManaging, fromManaged)
	go managing.Start(fromManaging, fromManaged)

	_, _ = fmt.Scanln()
}
