package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"time"
)

func generateSinInput(frequency float64, amplitude float64, duration time.Duration) []float64 {
	samplingRate := 44100 // Adjust this according to your requirements
	numSamples := int(duration.Seconds() * float64(samplingRate))
	sinInput := make([]float64, numSamples)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(samplingRate)
		sinInput[i] = amplitude * math.Sin(2*math.Pi*frequency*t)
	}

	return sinInput
}

func main() {
	frequency := 10.0       // Frequency of the sine wave in Hz (A4 note)
	amplitude := 100.0      // Amplitude of the sine wave
	duration := time.Second // Duration of the signal in seconds

	sinInput := generateSinInput(frequency, amplitude, duration)

	// Print or use sinInput as needed
	//fmt.Println(sinInput)

	for i := range sinInput {
		fmt.Println(sinInput[i])
	}
	// If you want to visualize the waveform, you can use a plotting library or external tool
	// For example, you can use gnuplot by writing the values to a temporary file
	//writeToTempFile(sinInput)
}

func writeToTempFile(data []float64) {
	file, err := os.CreateTemp("", "sin_wave_data.txt")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return
	}
	defer file.Close()

	for _, value := range data {
		_, err := fmt.Fprintln(file, value)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	// Open the file using gnuplot for visualization
	cmd := exec.Command("gnuplot", "-p", "-e", fmt.Sprintf(`plot "%s" with lines`, file.Name()))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running gnuplot:", err)
	}
}
