package main

import (
	"bufio"
	"fmt"
	"main.go/shared"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Data struct {
	QueueLength int
	PC          int
	Rate        float64
	Goal        float64
}

func main() {

	//et := flag.String("execution-type", "", "execution-type is a string")
	//inf := flag.String("input-file", "", "input-file is a string")
	//outf := flag.String("output-file", "", "output-file is a string")
	//flag.Parse()
	outf := "Experiment-PIWithTwoDegreesOfFreedom-Ziegler"
	b := "Experiment-PIWithTwoDegreesOfFreedom-Ziegler-"
	inputFiles := []string{}

	for i := 1; i <= 10; i++ {
		fileName := b + strconv.Itoa(i)
		inputFiles = append(inputFiles, fileName)
	}

	for i := 0; i < len(inputFiles); i++ {
		et := shared.Experiment
		//calcStatistics(*et, *inf, *outf)
		calcStatistics(et, inputFiles[i], outf)
	}
}

func calcStatistics(et, inf, outf string) {
	var includeHeader bool
	var outFile *os.File
	var err error

	// check if output file exist to place/not place file header
	filePath := shared.DataDir + "\\" + outf + ".csv"
	if shared.FileExists(filePath) {
		// open output file
		outFile, err = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		defer outFile.Close()
		includeHeader = false
	} else {
		// create output file
		outFile, err = os.Create(filePath)
		includeHeader = true
	}

	// generate data - raw files
	data := readFile(inf)

	if et == shared.Experiment {
		if includeHeader {
			fmt.Fprintf(outFile, "File;Tunning;RMSE;NMRSE;MAE;MAPE;SMAPE;R2;ITAE;ISE;Control Effort;CC;Goal Range \n")
		}
		fmt.Fprintf(outFile, "%v.csv;%.6f;%.6f;%.6f;%.6f;%.6f;%.6f;%.6f;%.6f;%.6f;%.6f;%.6f\n", inf, rmse(data), nmrse(data), mae(data), mape(data), smape(data), r2(data), itae(data), ise(data), controlEffort(data), cc(data), goalRange(data))
	} else {
		means := calculateMeans(data)

		// order by keys
		keys := make([]string, 0, len(means))
		for k := range means {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fmt.Fprintf(outFile, "%v;%v;%.6f;%v\n", 0, k, means[k], 0)
		}
	}
}

func calculateMeans(d []Data) map[string]float64 {
	n := len(d)
	sumLevel := make(map[int]float64)
	sizeLevel := make(map[int]int)
	means := make(map[string]float64)
	for i := 0; i < n; i++ {
		sumLevel[d[i].PC] = sumLevel[d[i].PC] + d[i].Rate
		sizeLevel[d[i].PC]++
	}
	for k, v := range sumLevel {
		key := strconv.Itoa(k)
		means[shared.PadLeft(key, 5)] = float64(v) / float64(sizeLevel[k])
	}
	return means
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
	filePath := shared.DataDir + "\\" + name + ".csv"
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
