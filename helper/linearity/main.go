package main

import (
	"fmt"
	"github.com/montanaflynn/stats"
)

func main() {

	x := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	y := []float64{150.8, 292.2, 543.0, 669.4, 840.8, 1011.0, 1153.6, 1273.4, 2087.2} // Replace with your actual data

	correlation, err := stats.Correlation(x, y)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(correlation)
}
