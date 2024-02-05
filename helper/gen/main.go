package main

import (
	"flag"
	"fmt"
	"main.go/shared"
	"os"
	"strings"
	"time"
)

func main() {

	// execution type
	et := flag.String("execution-type", "", "execution-type is a string")
	ct := flag.String("controller-type", "", "controller-type is a string")
	t := flag.String("tunning", "", "tunning is a string")
	f := flag.String("output-file", "", "output-file is a string")
	flag.Parse()

	if flag.NFlag() != 4 {
		shared.ErrorHandler(shared.GetFunction(), "Four parameters are necessary: execution-type, controller-type, tunning, output-file")
	}

	// configure dockerfile dir
	shared.ConfigureDockerfileDir(shared.DockerfilesDir)

	dockerFile := map[string]string{}
	dockerFile = configureDockerFile(*et, *ct, *t, *f)

	for k, v := range dockerFile {
		//fmt.Println(k, v)

		// create docker file
		df, err := os.Create(shared.DockerfilesDir + "\\" + k)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		defer df.Close()

		// write to docker file
		fmt.Fprintf(df, v)
	}
	// create batch file
	createBatFile(dockerFile)
}

func configureDockerFile(et, ct, t, f string) map[string]string {

	// header - common to all executions
	h := "# This file has been generated automatically at " + time.Now().String() + "\n" +
		"FROM golang:1.19\n" +
		"WORKDIR /app\n" +
		"COPY go.mod go.sum ./\n" +
		"RUN go mod download \n" +
		"COPY ./ ./ \n" +
		"RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/cli.go\n" +
		"ENV GOROOT=/usr/local/go/bin/\n"

	// command general parameters
	params := map[string]string{}
	params["execution-type"] = et
	params["controller-type"] = ct
	params["tunning"] = t
	params["output-file"] = f
	params["min"] = shared.MinPC
	params["max"] = shared.MaxPC
	params["monitor-interval"] = shared.MonitorInterval
	params["prefetch-count"] = shared.InitialPC
	params["set-point"] = shared.SetPoint
	params["direction"] = shared.Direction
	params["dead-zone"] = shared.DeadZone
	params["hysteresis-band"] = shared.HysteresisBand
	params["beta"] = shared.Beta
	params["alfa"] = shared.Alfa

	key := ""
	// pid parameters
	if ct == shared.SmoothingPid ||
		ct == shared.IncrementalFormPid ||
		ct == shared.ErrorSquarePidFull ||
		ct == shared.ErrorSquarePidProportional ||
		ct == shared.DeadZonePid ||
		ct == shared.GainScheduling ||
		ct == shared.PIwithTwoDegreesOfFreedom ||
		ct == shared.WindUp ||
		ct == shared.SetpointWeighting {
		key = shared.BasicPid + t
	} else {
		key = ct + t
	}
	v, ok := shared.Kp[key]
	if ok {
		params["kp"] = v
	} else {
		params["kp"] = "1.0"
	}
	v, ok = shared.Ki[key]
	if ok {
		params["ki"] = v
	} else {
		params["ki"] = "1.0"
	}
	v, ok = shared.Kd[key]
	if ok {
		params["kd"] = v
	} else {
		params["kd"] = "1.0"
	}

	r := map[string]string{}
	fileName := "Dockerfile" + "-" + params["execution-type"] + "-" + params["controller-type"] + "-" + params["tunning"]
	content := h
	p := ""
	for k, v := range params {
		p += "\"" + "-" + k + "=" + v + "\","
	}
	// remove last ´,´
	cmd := "CMD [\"./subscriber\"," + p + "]"
	cmd = strings.Replace(cmd, ",]", "]", 100)
	content += cmd
	r[fileName] = content

	return r
}

func createBatFile(list map[string]string) {
	// define list of docker files include in bat file
	listCommand := "set list="
	for k, _ := range list {
		listCommand += k + " "
	}
	listCommand += "\n"

	batFileContent := "rem This file has been generated automatically at " + time.Now().String() + "\n" +
		"@echo off \n" +
		//	"docker stop some-rabbit \n" +
		//	"docker rm some-rabbit\n" +
		//	"docker run -d --memory=\"6g\" --cpus=\"5.0\" --name some-rabbit -p 5672:5672 rabbitmq\n" +
		"timeout /t 10\n" +
		"docker stop publisher\n" +
		"docker rm publisher\n" +
		"docker stop subscriber\n" +
		"docker rm subscriber\n"

	batFileContent += listCommand + "\n"
	batFileContent +=
		"for %%x in (%list%) do (\n" +
			"echo %%x \n" +
			"   copy " + shared.DockerfilesDir + "\\" + "%%x Dockerfile\n" +
			"   docker build --tag subscriber .\n" +
			"   docker run --rm --name some-subscriber --memory=\"1g\" --cpus=\"1.0\" -v " + shared.DataDir + ":" + shared.DockerDir + " subscriber\n" +
			"   del " + shared.DockerfilesDir + "\\" + "%%x \n" +
			"   echo y| docker volume prune \n" +
			"   echo y| docker image prune \n" +
			")\n"
	//"docker stop some-rabbit \n" +
	//"docker rm some-rabbit\n"

	// create batch file
	fileName := shared.BatchFileExperiments
	bf, err := os.Create(shared.BatchfilesDir + "\\" + fileName)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	defer bf.Close()

	fmt.Fprintf(bf, "%v", batFileContent)
}
