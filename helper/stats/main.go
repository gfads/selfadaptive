package main

import (
	"bufio"
	"fmt"
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

type Metrics struct {
	RMSE  float64
	NRMSE float64
	MAPE  float64
	SMAPE float64
	ITAE  float64
	IAE   float64
	ISE   float64
	CE    float64
	R2    float64
	GR    float64
}

func main() {
	allData := map[string]Metrics{}
	/*controllers := []string{"hpa-fixed", "mypi-fixed", "mypid-fixed",
	"pitf10faster-fixed", "pidtf10faster-fixed", "pitf21faster-fixed", "pidtf21faster-fixed",
	"piziegler-fixed", "pidziegler-fixed", "picohen-fixed", "pidcohen-fixed", "piamigo-fixed",
	"pidamigo-fixed", "gain-fixed", "astar-fixed",
	"piddeadzone-fixed","pidsmoothing-fized","pidincrementalform-fixed"}
	*/
	controllers := []string{"hpa-variable", "mypi-variable", "mypid-variable",
		"pitf10faster-variable", "pidtf10faster-variable", "pitf21faster-variable", "pidtf21faster-variable",
		"piziegler-variable", "pidziegler-variable", "picohen-variable", "pidcohen-variable", "piamigo-variable",
		"pidamigo-variable", "gain-variable", "astar-variable",
		"piddeadzone-variable", "pidsmoothing-variable", "pidincrementalform-variable"} //outFile := "all-variable-summary.csv"
	outFile := "all-fixed-summary.csv"

	for c := 0; c < len(controllers); c++ {
		data := []Data{}
		fileName := controllers[c]
		d := readFile(fileName)
		for j := 0; j < len(d); j++ {
			data = append(data, d[j])
		}
		// calculate metrics
		allData[controllers[c]] = calcMetrics(data)
	}
	saveMetrics(outFile, allData)
	fmt.Println("Stats Finished!!!")
}

func calcMetrics(d []Data) Metrics {
	r := Metrics{}
	r.RMSE = rmse(d)
	r.NRMSE = nmrse(d)
	r.MAPE = mape(d)
	r.SMAPE = smape(d)
	r.ITAE = itae(d)
	r.IAE = iae(d)
	r.ISE = ise(d)
	r.CE = controlEffort(d)
	r.R2 = r2(d)
	r.GR = goalRange(d)

	return r
}

func saveMetrics(fileName string, allData map[string]Metrics) {
	var outFile *os.File
	var err error

	filePath := shared.DataDir + "\\" + fileName
	outFile, err = os.Create(filePath)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	fmt.Fprintf(outFile, "Controller;RMSE;NMRSE;MAPE;SMAPE;ITAE;IAE;ISE;Control Effort;R2;Goal Range \n")
	for k := range allData {
		fmt.Fprintf(outFile, "%v;%.6f;%.6f;%.6f;%.6f;%.6f;%.6f;%.6f;%.6f;%.6f;%.6f\n", k,
			allData[k].RMSE, allData[k].NRMSE, allData[k].MAPE, allData[k].SMAPE,
			allData[k].ITAE, allData[k].IAE, allData[k].ISE, allData[k].CE,
			allData[k].R2, allData[k].GR)
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
		r += float64(d[i].PC)
	}
	return r
}

func r2(d []Data) float64 {
	n := len(d)
	sXY := 0.0
	sX := 0.0
	sY := 0.0
	sX2 := 0.0
	sY2 := 0.0

	for i := 0; i < n; i++ {
		sXY += d[i].Rate * d[i].Goal
		sX += d[i].Rate
		sY += d[i].Goal
		sX2 += math.Pow(d[i].Rate, 2.0)
		sY2 += math.Pow(d[i].Goal, 2.0)
	}

	//fmt.Printf("%.10f %.10f %.10f\n ", float64(n)*sXY, sX*sY, float64(n))
	num := sXY*float64(n) - sX*sY
	exp1 := float64(n)*sX2 - sX*sX
	exp2 := float64(n)*sY2 - sY*sY
	den := math.Sqrt(exp1 * exp2)

	r := math.Pow(num/den, 2.0)

	return r
}

func smape(d []Data) float64 {
	s := 0.0
	n := len(d)
	for i := 0; i < n; i++ {
		s += 2 * math.Abs(d[i].Rate-d[i].Goal) / (d[i].Rate + d[i].Goal)
	}
	r := s / float64(n) * 100.0
	return r
}

func itae(d []Data) float64 {
	r := 0.0
	n := len(d)
	t, _ := strconv.Atoi(shared.MonitorInterval)
	for i := 0; i < n; i++ {
		r += math.Abs(d[i].Rate-d[i].Goal) * float64(t)
		temp, _ := strconv.Atoi(shared.MonitorInterval)
		t += temp
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

func iae(d []Data) float64 {
	r := 0.0
	n := len(d)
	for i := 0; i < n; i++ {
		r += math.Abs(d[i].Rate - d[i].Goal)
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
