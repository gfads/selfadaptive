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
	Behaviours map[string]func()
}

func NewManagedElement() *ManagedElement {

	temp := make(map[string]func())
	r := ManagedElement{Behaviours: temp}

	r.Behaviours["DefaultBehaviour"] = r.PlainText
	r.Behaviours["WeakCryptography"] = r.WeakCryptography
	r.Behaviours["MediumCryptography"] = r.MediumCryptography
	r.Behaviours["StrongCryptography"] = r.StrongCryptography

	return &r
}

func (m ManagedElement) Run(toManaging chan map[string]func(), fromManaging chan shared.ToManagedChan) { // Business logic

	m.PlainText()

	for {
		toManaging <- m.Behaviours // To managing
		b := <-fromManaging        // From managing

		behaviour := b.Behaviours[b.SelectedBehaviour] // new behaviour
		behaviour()                                    // change to new behaviour
	}
}

func (ManagedElement) PlainText() { // plain text
	fmt.Printf("[Default] -> [Plain Text] '%s'\n", shared.PlainText)
}

func (ManagedElement) WeakCryptography() {
	fmt.Printf("[Weak] '%s'\n", shared.EncryptMessage(shared.PlainText, shared.Keys32[0]))
}

func (ManagedElement) MediumCryptography() {
	fmt.Printf("[Medium] '%s'\n", shared.EncryptMessage(shared.PlainText, shared.Keys32[1]))
}

func (ManagedElement) StrongCryptography() {
	fmt.Printf("[Strong] '%s'\n", shared.EncryptMessage(shared.PlainText, shared.Keys32[2]))
}
