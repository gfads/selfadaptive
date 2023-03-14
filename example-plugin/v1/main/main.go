package main

import (
	"fmt"
	"selfadaptive/example-plugin/v1/env"
	"selfadaptive/example-plugin/v1/mnged"
	"selfadaptive/example-plugin/v1/mnging"
	"selfadaptive/shared/channeltypes"
)

func main() {

	fromManaged := make(chan channeltypes.TypeChanManaging)
	fromManaging := make(chan []func())

	managed := mnged.NewManagedElement()
	managing := mnging.NewManagingSystem()
	e := env.NewEnvironment()

	go e.Start()
	go managed.Start(fromManaging, fromManaged)
	go managing.Start(fromManaging, fromManaged)

	fmt.Scanln()
}
