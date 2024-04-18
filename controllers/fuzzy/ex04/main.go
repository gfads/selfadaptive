package main

/*
func main() {
	// Create a fuzzy controller
	c := fuzzy.NewController()

	// Define fuzzy variables: temperature and fan speed
	temperature := c.NewInput("temperature")
	temperature.SetRange(0, 100) // Celsius

	speed := c.NewOutput("speed")
	speed.SetRange(0, 100) // Percentage of maximum speed

	// Define membership functions for temperature: low, medium, high
	temperature.Triangular("low", 0, 20, 40)
	temperature.Triangular("medium", 30, 50, 70)
	temperature.Triangular("high", 60, 80, 100)

	// Define membership functions for fan speed: low, medium, high
	speed.Triangular("low", 0, 25, 50)
	speed.Triangular("medium", 30, 50, 70)
	speed.Triangular("high", 60, 75, 100)

	// Define fuzzy rules
	c.Rule("if temperature is low then speed is high")
	c.Rule("if temperature is medium then speed is medium")
	c.Rule("if temperature is high then speed is low")

	// Generate a random temperature value
	rand.Seed(time.Now().UnixNano())
	temperatureValue := rand.Float64() * 100

	// Calculate fan speed based on temperature
	result := c.Run(map[string]float64{"temperature": temperatureValue})

	// Get the inferred speed
	fanSpeed := result["speed"]

	// Output the result
	fmt.Printf("Temperature: %.2f°C\n", temperatureValue)
	fmt.Printf("Fan Speed: %.2f%%\n", fanSpeed)
}
*/
