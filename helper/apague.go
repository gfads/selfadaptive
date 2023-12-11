package main

import (
	"bufio"
	"fmt"
	"main.go/rabbitmq/adaptationlogic"
	"main.go/shared"
	"os"
	"strconv"
	"strings"
)

func main() {

	data := []adaptationlogic.AdjustmenstInfo{}
	data = append(data, adaptationlogic.AdjustmenstInfo{PC: 1, Rate: 0.0})
	c := readContent("C:\\Users\\user\\go\\selfadaptive\\rabbitmq\\data\\journal\\sin-training-data-200.csv")

	d := generateDataInfo(c)
	t := adaptationlogic.TrainingInfo{Data: d, TypeName: shared.BasicPid}
	t1 := adaptationlogic.CalculateRootLocusGains(t)
	fmt.Println(t1.Kp)
	fmt.Println(t1.Ki)
	fmt.Println(t1.Kd)
}

func generateDataInfo(c []string) []adaptationlogic.AdjustmenstInfo {
	r := []adaptationlogic.AdjustmenstInfo{}

	for i := range c {
		l := c[i]
		d := strings.Split(l, ";")
		pc, _ := strconv.Atoi(d[1])
		if s, err := strconv.ParseFloat(d[2], 64); err == nil {
			r = append(r, adaptationlogic.AdjustmenstInfo{PC: pc, Rate: s})
		}
	}
	return r
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
