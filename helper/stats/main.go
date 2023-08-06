package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
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

const nameFilter = "raw-sin-" // TODO

func main() {

	// open output stat file
	statFile, err := os.Create(shared.DataDir + "\\" + shared.StatiticsFileName)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	defer statFile.Close()

	// read folder of files
	files, err := ioutil.ReadDir(shared.DataDir)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	for i := range files {
		fmt.Println("HERE", files[i].Name())
	}
	// generate data
	fmt.Fprintf(statFile, "Controller;Tunning;RMSE;NMRSE;MAE;MAPE;SMAPE;R2;ITAE;ISE;Control Effort;CC;Goal Range \n")
	for f := range files {
		if strings.Contains(files[f].Name(), nameFilter) {
			data := readFile(files[f].Name())
			//i1 := strings.Index(files[f].Name(), ".csv")
			//temp := strings.Split(files[f].Name()[len(nameFilter):i1], "-")
			//controller := temp[0]
			//tunning := temp[1]
			fmt.Fprintf(statFile, "%v;%.6f;%.6f;%.6f;%.6f;%.6f;%.6f;%.6f;%.6f;%.6f;%.6f;%.6f\n", files[f].Name(), rmse(data), nmrse(data), mae(data), mape(data), smape(data), r2(data), itae(data), ise(data), controlEffort(data), cc(data), goalRange(data))
		}
	}
}

func goalRange(d []Data) float64 {
	n := len(d)
	r := 0
	for i := 0; i < n; i++ {
		if d[i].Rate > d[i].Goal*0.8 && d[i].Rate < d[i].Goal*1.2 {
			r++
		}
	}
	return float64(r) / float64(n) * 100.0
}

func cc(d []Data) float64 {
	n := len(d)

	// calculate mean
	x := 0.0
	y := 0.0
	for i := 0; i < n; i++ {
		x += d[i].Rate
		y += float64(d[i].PC)
	}
	mX := x / float64(n)
	mY := y / float64(n)
	numerator := 0.0
	dX := 0.0
	dY := 0.0
	for i := 0; i < n; i++ {
		numerator += (d[i].Rate - mX) * (float64(d[i].PC) - mY)
		dX += math.Pow(d[i].Rate-mX, 2.0)
		dY += math.Pow(float64(d[i].PC)-mY, 2.0)
	}
	d1 := math.Sqrt(dX)
	d2 := math.Sqrt(dY)

	r := numerator / d1 * d2

	return r
}

func controlEffort(d []Data) float64 {
	n := len(d)
	r := 0.0

	// calculate mean
	for i := 0; i < n; i++ {
		r += math.Pow(float64(d[i].PC), 2.0)
	}
	return r
}

func r2(d []Data) float64 {
	tss := 0.0
	rss := 0.0
	n := len(d)
	temp := 0.0

	// calculate mean
	for i := 0; i < n; i++ {
		temp += d[i].Rate
	}
	mean := temp / float64(n)
	for i := 0; i < n; i++ {
		rss += math.Pow(d[i].Rate-d[i].Goal, 2.0)
		tss += math.Pow(d[i].Rate-mean, 2.0)
	}
	r := 1 - (rss / tss)
	return r
}

func smape(d []Data) float64 {
	s := 0.0
	n := len(d)
	for i := 0; i < n; i++ {
		s += math.Abs(d[i].Rate-d[i].Goal) / ((d[i].Rate + d[i].Goal) / 2.0)
	}
	r := s / float64(n) * 100.0
	return r
}

func itae(d []Data) float64 {
	r := 0.0
	n := len(d)
	for i := 0; i < n; i++ {
		r += math.Abs(d[i].Rate - d[i].Goal)
	}
	return r
}

func ise(d []Data) float64 {
	r := 0.0
	n := len(d)
	for i := 0; i < n; i++ {
		r += math.Pow(math.Abs(d[i].Rate-d[i].Goal), 2.0)
	}
	return r
}

func mape(d []Data) float64 {
	s := 0.0
	n := len(d)
	for i := 0; i < n; i++ {
		s += math.Abs(d[i].Rate-d[i].Goal) / d[i].Rate
	}
	r := s / float64(n) * 100.0
	return r
}

func mae(d []Data) float64 {
	s := 0.0
	n := len(d)
	for i := 0; i < n; i++ {
		s += math.Abs(d[i].Rate - d[i].Goal)
	}
	r := s / float64(n)
	return r
}

func rmse(d []Data) float64 {
	s := 0.0
	n := len(d)
	for i := 0; i < n; i++ {
		s += math.Pow(d[i].Rate-d[i].Goal, 2.0)
	}
	r := math.Sqrt(s / float64(n))

	return r
}

func nmrse(d []Data) float64 {
	n := len(d)

	// find max and min
	max := 0.0
	min := 10000000.0
	for i := 0; i < n; i++ {
		if d[i].Rate < min {
			min = d[i].Rate
		}
		if d[i].Rate > max {
			max = d[i].Rate
		}
	}
	r := rmse(d) / (max - min)

	return r
}
func readFile(name string) []Data {

	data := []Data{}

	filePath := shared.DataDir + "/" + name

	readFile, err := os.Open(filePath)

	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		l := fileScanner.Text()

		// check format of number
		if strings.Contains(l, ",") {
			l = strings.ReplaceAll(l, ",", ".")
		}
		s := strings.Split(l, ";")

		// queue length
		ql, err := strconv.Atoi(s[0])
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		// pc
		pc, err := strconv.Atoi(s[1])
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
