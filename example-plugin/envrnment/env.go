package envrnment

import (
	"fmt"
	"math/rand"
	"selfadaptive/shared"
	"time"
)

type Environment struct {
	SecurityLevel string
	Plugins       []func()
}

func NewEnvironment() *Environment {
	e := Environment{SecurityLevel: shared.Secure, Plugins: shared.LoadPlugins(shared.SourcesDir, shared.ExecutablesDir)}
	return &e
}

func (e *Environment) Run() {
	for {
		// generate a new security level randomly
		time.Sleep(10 * time.Second)
		rand.Seed(time.Now().UnixNano())
		sls := shared.EnvironmentSecurityLevels
		sl := sls[rand.Intn(len(sls))]

		// update security level
		e.SecurityLevel = sl
	}
}

func (e Environment) Sense() (string, []func()) {
	r1 := e.SecurityLevel
	r2 := shared.LoadPlugins(shared.SourcesDir, shared.ExecutablesDir)

	fmt.Printf("[%s] -> ", r1)

	return r1, r2
}
