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

func main() {
	dir := "C:\\Users\\user\\go\\selfadaptive\\rabbitmq\\data"
	filePattern := "ZieglerTraining-BasicP-Ziegler-"
	t := "AMIGO" // ´´ - RooLocus, ´Ziegler´, ´Cohen´, ´AMIGO´

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	allContents := map[string][]string{}
	allKps := map[string]float64{}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), filePattern) {
			path := dir + "\\" + file.Name()
			allContents[path] = readContent(path)
		}
	}

	s := 0.0
	for k, _ := range allContents {
		kp := extractKp(allContents[k], t)
		allKps[k] = kp
		s += kp
	}
	mean := s / float64(len(allKps))
	gains := findBestMatch(allContents, allKps, t, mean)
	fmt.Println(gains)
}

func findBestMatch(allContents map[string][]string, allKps map[string]float64, t string, mean float64) string {
	f := ""
	r := ""
	min := 10000000.0

	for k, v := range allKps {
		n := math.Abs(math.Abs(v) - math.Abs(mean))
		if n < min {
			min = n
			f = k
		}
	}
	for l := range allContents[f] {
		if strings.Contains(allContents[f][l], "-kp=") && strings.Contains(allContents[f][l], t) {
			r = f + allContents[f][l]
		}
	}
	return r
}
func extractKp(f []string, t string) float64 {
	kp := 0.0
	var err error
	for i := range f {
		if strings.Contains(f[i], "-kp=") && strings.Contains(f[i], t) {
			x := strings.Split(f[i], "\",")
			j := strings.Index(x[0], "-kp=")
			kp, err = strconv.ParseFloat(x[0][j+4:], 64)
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), "Erro na conversão")
			}
		}
	}
	return kp
}
func readContent(p string) []string {
	readFile, err := os.Open(p)

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}
	readFile.Close()

	return fileLines
}
