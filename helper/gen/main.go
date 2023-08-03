package main

import (
	"fmt"
	"main.go/shared"
	"os"
	"time"
)

func main() {
	controllers := shared.ControllerTypes
	tunnings := shared.TunningTypes
	dockerFileNames := []string{}
	pExtra := loadExtraParameters(controllers, tunnings)

	// create Dockerfiles
	for c := 0; c < len(controllers); c++ {
		for t := 0; t < len(tunnings); t++ {
			dockerFileNames = append(dockerFileNames, createDockerFile(controllers[c], tunnings[t], pExtra))
		}
	}

	// create execute
	createBat(dockerFileNames)
}

func loadExtraParameters(c, t []string) map[string]string {

	r := make(map[string]string)

	// Non PID controllers
	r[shared.HPA+shared.None] = ""
	r[shared.BasicOnoff+shared.None] = ""
	r[shared.DeadZoneOnoff+shared.None] = "\"-dead-zone=" + shared.DeadZone + "\""
	r[shared.HysteresisOnoff+shared.None] = "\"-hysteresis-band=" + shared.HysteresisBand + "\""
	r[shared.AsTAR+shared.None] = "\"-hysteresis-band=" + shared.HysteresisBand + "\""

	// PID controllers
	for i := 0; i < len(c); i++ {
		for j := 0; j < len(t); j++ {
			key := c[i] + t[j]
			keyGain := ""
			if c[i] == shared.BasicP || c[i] == shared.BasicPi {
				keyGain = c[i] + t[j]
			} else {
				keyGain = shared.BasicPid + t[j] // use the same key whatever the PID controller (see shared)
			}
			r[key] =
				"\"-kp=" + shared.Kp[keyGain] + "\", " +
					"\"-ki=" + shared.Ki[keyGain] + "\", " +
					"\"-kd=" + shared.Kd[keyGain] + "\""
		}
	}

	// Extra parameters of PID controllers
	r[shared.DeadZonePid+shared.RootLocus] += ",\"-dead-zone=" + shared.DeadZone + "\""
	r[shared.DeadZonePid+shared.Ziegler] += ",\"-dead-zone=" + shared.DeadZone + "\""
	r[shared.DeadZonePid+shared.Cohen] += ",\"-dead-zone=" + shared.DeadZone + "\""
	r[shared.DeadZonePid+shared.Amigo] += ",\"-dead-zone=" + shared.DeadZone + "\""

	r[shared.PIwithTwoDegreesOfFreedom+shared.RootLocus] = ",\"-beta=+" + shared.Beta + "\""
	r[shared.PIwithTwoDegreesOfFreedom+shared.Ziegler] = ",\"-beta=+" + shared.Beta + "\""
	r[shared.PIwithTwoDegreesOfFreedom+shared.Cohen] = ",\"-beta=+" + shared.Beta + "\""
	r[shared.PIwithTwoDegreesOfFreedom+shared.Amigo] = ",\"-beta=+" + shared.Beta + "\""
	r[shared.SetpointWeighting+shared.RootLocus] = ",\"-alfa=" + shared.Alfa + "," + shared.Beta + "\""
	r[shared.SetpointWeighting+shared.Ziegler] = ",\"-alfa=" + shared.Alfa + "," + shared.Beta + "\""
	r[shared.SetpointWeighting+shared.Cohen] = ",\"-alfa=" + shared.Alfa + "," + shared.Beta + "\""
	r[shared.SetpointWeighting+shared.Amigo] = ",\"-alfa=" + shared.Alfa + "," + shared.Beta + "\""

	return r
}

func createDockerFile(c, t string, pExtra map[string]string) string {

	// remove old docker files
	err := os.RemoveAll(shared.DockerfilesDir)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// recreate docker files folder
	err = os.MkdirAll(shared.DockerfilesDir, 0750)
	if err != nil && !os.IsExist(err) {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// set tunning of non pid controllers to None
	if t == shared.AsTAR || t == shared.HPA || t == shared.BasicOnoff ||
		t == shared.DeadZoneOnoff || t == shared.HysteresisOnoff {
		t = shared.None
	}

	// configure general "command" parameters
	min := "\"-min=" + shared.MinPC + "\""
	max := "\"-max=" + shared.MaxPC + "\""
	monitorInterval := "\"-monitor-interval=" + shared.MonitorInterval + "\""
	executionType := "\"-execution-type=" + shared.Experiment + "\""
	adaptability := "\"-is-adaptive=" + shared.Adaptability + "\""
	prefetchCount := "\"-prefetch-count=" + shared.PrefetchCountInitial + "\""
	setPoint := "\"-set-point=" + shared.SetPoint + "\""
	direction := "\"-direction=" + shared.Direction + "\""
	controller := "\"-controller-type=" + c + "\""
	tunning := "\"-tunning=" + t + "\""

	// create docker file
	fileName := "Dockerfile-" + c + "-" + t
	df, err := os.Create(shared.DockerfilesDir + "\\" + fileName)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	defer df.Close()

	// docker file content
	dockerFileContent := "# Self generated file at " + time.Now().String() + "\n" +
		"FROM golang:1.19\n" +
		"WORKDIR /app\n" +
		"COPY go.mod go.sum ./\n" +
		"RUN go mod download \n" +
		"COPY ./ ./ \n" +
		"RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go\n" +
		"ENV GOROOT=/usr/local/go/bin/\n" +
		"CMD [\"./subscriber\"," + controller + "," + tunning + "," + adaptability + "," + executionType + "," + monitorInterval + "," + prefetchCount + "," + max + "," + min + "," + setPoint + "," + direction + "," + pExtra[c+t] + "]"

	// write in docker file
	fmt.Fprintf(df, dockerFileContent)

	return fileName
}

func createBat(list []string) {
	// define list of docker files include in bat file
	listCommand := "set list="
	for i := 0; i < len(list); i++ {
		listCommand += list[i] + " "
	}
	listCommand += "\n"

	batFileContent := "rem Self generated file at " + time.Now().String() + "\n" +
		"@echo off \n" +
		"docker stop some-rabbit \n" +
		"docker rm some-rabbit\n" +
		"docker run -d --memory=\"6g\" --cpus=\"5.0\" --name some-rabbit -p 5672:5672 rabbitmq\n" +
		"timeout /t 20\n" +
		"echo Removing previous containers\n" +
		"docker stop publisher\n" +
		"docker rm publisher\n" +
		"docker stop subscriber\n" +
		"docker rm subscriber\n" +
		"echo Removing previous volumes\n" +
		"echo y | docker volume prune\n"

	batFileContent += listCommand + "\n"
	batFileContent +=
		"echo ****** BEGIN OF EXPERIMENTS *******\n" +
			"for %%x in (%list%) do (\n" +
			"echo %%x \n" +
			"   copy " + shared.DockerfilesDir + "\\" + "%%x Dockerfile\n" +
			"   docker build --tag subscriber .\n" +
			"   docker run --rm --name some-subscriber --memory=\"1g\" --cpus=\"1.0\" -v " + shared.DataDir + ":" + shared.DockerDir + " subscriber\n" +
			"   del " + shared.DockerfilesDir + "\\" + "%%x \n" +
			")\n" +
			"echo ****** END OF EXPERIMENTS *******\n"

	// create batch file
	fileName := shared.BatchFileExperiments
	bf, err := os.Create(shared.BatchfilesDir + "\\" + fileName)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	defer bf.Close()

	fmt.Fprintf(bf, batFileContent)
}
