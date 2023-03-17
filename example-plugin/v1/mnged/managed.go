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

	r.Behaviours = append(r.Behaviours, r.behaviour00)
	r.Behaviours = append(r.Behaviours, r.behaviour01)
	r.Behaviours = append(r.Behaviours, r.behaviour02)

	return &r
}

func (m ManagedElement) Start(toManaging chan []func(), fromManaging chan shared.TypeChanManaging) { // Business logic
	for {
		toManaging <- m.Behaviours // To managing
		b := <-fromManaging        // From managing

		behaviour := b.Functions[b.N] // new behaviour
		behaviour()                   // change to new behaviour
	}
}

func (m ManagedElement) behaviour00() {
	for i := 0; i < 100; i++ {
		fmt.Print(string(shared.ColorBehaviours[0]), "+")
	}
	fmt.Print(shared.ColorReset)
}
func (m ManagedElement) behaviour01() {
	for i := 0; i < 100; i++ {
		fmt.Print(string(shared.ColorBehaviours[1]), "+")
	}
	fmt.Print(shared.ColorReset)
}
func (m ManagedElement) behaviour02() {
	for i := 0; i < 100; i++ {
		fmt.Print(string(shared.ColorBehaviours[2]), "+")
	}
	fmt.Print(shared.ColorReset)
}
