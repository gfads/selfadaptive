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
const FileNameRadical = "raw-sin-3-"

func main() {

	// read files
	files, err := ioutil.ReadDir(DataDir)
	if err != nil {
		log.Fatal(err)
	}

	for f := range files {
		if strings.Contains(files[f].Name(), FileNameRadical) {
			data := readFile(files[f].Name())
			i1 := strings.Index(files[f].Name(), ".csv")
			temp := strings.Split(files[f].Name()[len(FileNameRadical):i1], "-")
			controller := temp[0]
			tunning := temp[1]
			fmt.Printf("%v;%v;%.6f;%.6f;%.6f;%.6f;%.6f\n", controller, tunning, rmse(data), mae(data), mape(data), smape(data), r2(data))
		}
	}
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
		tss += math.Pow(d[i].Rate-mean, 2.0)
		rss += math.Pow(d[i].Rate-d[i].Goal, 2.0)
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
