package fuzzy

import (
	"fmt"
	"main.go/controllers/def/info"
	"main.go/shared"
	"os"
)

type Controller struct {
	Info info.Controller
}

func (c *Controller) Initialise(p ...float64) {

	if len(p) < 4 {
		fmt.Printf("Error: '%s controller requires 4 info (input.center, input.width, output.center, output.width) \n", shared.Fuzzy)
		os.Exit(0)
	}

	c.Info.TypeName = shared.Fuzzy
	c.Info.InputSet.Center = p[0]
	c.Info.InputSet.Width = p[1]
	c.Info.OutputSet.Center = p[2]
	c.Info.OutputSet.Width = p[3]
}

// fuzzy rules
func (c Controller) fuzzyRule(input float64) float64 {
	// Rule 1: If input is low, print delay is low
	if input <= c.Info.InputSet.Center-c.Info.InputSet.Width {
		return c.Info.OutputSet.Center - c.Info.OutputSet.Width
	}
	// Rule 2: If input is medium, print delay is medium
	if input >= c.Info.InputSet.Center-c.Info.InputSet.Width && input <= c.Info.InputSet.Center+c.Info.InputSet.Width {
		return c.Info.OutputSet.Center
	}
	// Rule 3: If input is high, print delay is high
	if input >= c.Info.InputSet.Center+c.Info.InputSet.Width {
		return c.Info.OutputSet.Center + c.Info.OutputSet.Width
	}
	return 0 // Default
}

// Fuzzy inference
func (c Controller) fuzzyInference(input float64) float64 {
	// Apply fuzzy rules
	output := c.fuzzyRule(input)
	return output
}

func (c *Controller) Update(p ...float64) float64 {
	u := c.fuzzyInference(p[0])
	return u
}
