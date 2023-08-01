package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"main.go/shared"
	"math"
	"os"
	"strconv"
	"strings"
)

type Data struct {
	QueueLength int
	PC          int
	Rate        float64
	Goal        float64
}

const DataDir = "/Volumes/GoogleDrive/Meu Drive/go/selfadaptive/rabbitmq/data/"

func main() {

	// read files
	files, err := ioutil.ReadDir(DataDir)
	if err != nil {
		log.Fatal(err)
	}

	for f := range files {
		if strings.Contains(files[f].Name(), "raw-sin-3") {
			data := readFile(files[f].Name())
			fmt.Printf(files[f].Name()+" RMSE = %.6f\n", rmse(data))
		}
	}
}

func r2(d []Data) float64 {
	return 0
}

func mape(d []Data) float64 {
	return 0
}

func mae(d []Data) float64 {
	return 0
}

func rmse(d []Data) float64 {
	s := 0.0
	n := len(d)
	for i := 0; i < n; i++ {
		s += math.Pow(d[i].Rate-d[i].Goal, 2.0)
	}
	rmse := math.Sqrt(s / float64(n))

	return rmse
}

func readFile(name string) []Data {

	data := []Data{}

	filePath := DataDir + "/" + name

	readFile, err := os.Open(filePath)

	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		l := fileScanner.Text()
		s := strings.Split(l, ";")

		// queue length
		ql, err := strconv.Atoi(s[0])
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		// pc
		pc, err := strconv.Atoi(s[0])
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		// rate
		rate, err := strconv.ParseFloat(s[2], 32)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		// goal
		goal, err := strconv.ParseFloat(s[3], 32)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		d := Data{QueueLength: ql, PC: pc, Rate: rate, Goal: goal}
		data = append(data, d)
	}
	readFile.Close()

	return data
}
