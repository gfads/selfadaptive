package envrnment

import (
	"fmt"
	"math/rand"
	"selfadaptive/shared"
	"time"
)

type Environment struct {
	SecurityLevel int
	Plugins       []func()
}

func NewEnvironment() *Environment {
	e := Environment{SecurityLevel: shared.HighSecurityLevel, Plugins: shared.LoadPlugins(shared.SourcesDir, shared.ExecutablesDir)}
	return &e
}

func (e *Environment) Run() {
	for {
		e.SecurityLevel = shared.HighSecurityLevel
		time.Sleep(10 * time.Hour)
	}
}

func (Environment) Sense() (string, []func()) {
	r1 := shared.EnvironmentSecurityLevel[rand.Intn(2)]
	r2 := shared.LoadPlugins(shared.SourcesDir, shared.ExecutablesDir)

	fmt.Print("Env:: ", r1, " ")

	return r1, r2
}
