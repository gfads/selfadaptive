/*********************************************************************************
Author: Nelson S Rosa
Description: This code implements a simple managed system that prints something
on the screen. The managed system has three harded-code behaviours.
Date: 14/02/2023
*********************************************************************************/
package mnged

import (
	"fmt"
	"selfadaptive/shared"
)

type ManagedElement struct {
	Behaviours []func()
}

func NewManagedElement() *ManagedElement {

	r := ManagedElement{}

	r.Behaviours = append(r.Behaviours, r.defaultBehaviour)

	return &r
}

func (m ManagedElement) Run(toManaging chan []func(), fromManaging chan shared.ToManagedChan) { // Business logic

	m.defaultBehaviour()

	for {
		toManaging <- m.Behaviours // To managing
		b := <-fromManaging        // From managing

		behaviour := b.Behaviours[b.SelectedBehaviour] // new behaviour
		behaviour()                                    // change to new behaviour
	}
}

func (m ManagedElement) defaultBehaviour() {
	fmt.Println("Sent Message:", shared.PlainText)
}
