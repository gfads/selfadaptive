package main

import (
	"fmt"
	"main.go/shared"
	"os"
	"time"
)

func main() {
	dockerFileNames := []string{}
	pExtra := loadExtraParameters()

	// create Dockerfiles Non-PID Controllers
	for c := 0; c < len(shared.ControllerTypesNonPid); c++ {
		//	createDockerFile(shared.ControllerTypesNonPid[c], shared.None, pExtra, &dockerFileNames)
	}

	// create Dockerfiles PID Controllers
	for c := 0; c < len(shared.ControllerTypesPid); c++ {
		for t := 0; t < len(shared.TunningTypes); t++ {
			createDockerFile(shared.ControllerTypesPid[c], shared.TunningTypes[t], pExtra, &dockerFileNames)
		}
	}
	fmt.Println("Docker files created...")

	// create execute
	createBat(dockerFileNames)

	fmt.Println("Bat file created...")
}

func loadExtraParameters() map[string]string {

	r := make(map[string]string)

	// Non PID controllers
	r[shared.HPA+shared.None] = ""
	r[shared.BasicOnoff+shared.None] = ""
	r[shared.DeadZoneOnoff+shared.None] = "\"-dead-zone=" + shared.DeadZone + "\""
	r[shared.HysteresisOnoff+shared.None] = "\"-hysteresis-band=" + shared.HysteresisBand + "\""
	r[shared.AsTAR+shared.None] = "\"-hysteresis-band=" + shared.HysteresisBand + "\""

	// PID controllers
	for i := 0; i < len(shared.ControllerTypesPid); i++ {
		for j := 0; j < len(shared.TunningTypes); j++ {
			key := shared.ControllerTypesPid[i] + shared.TunningTypes[j]
			keyGain := shared.BasicPid + shared.TunningTypes[j] // use the same key whatever the PID controller (see shared)
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

func createDockerFile(c, t string, pExtra map[string]string, df *[]string) {

	// docker files folder
	basicDocker := "# Self generated file at " + time.Now().String() + "\n" +
		"FROM golang:1.19\n" +
		"WORKDIR /app\n" +
		"COPY go.mod go.sum ./\n" +
		"RUN go mod download \n" +
		"COPY ./ ./ \n" +
		"RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go\n" +
		"ENV GOROOT=/usr/local/go/bin/"

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
	fileName := "Dockerfile-" + c + "-" + t
	f, err := os.Create(shared.DockerfilesDir + "\\" + fileName)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	command := "CMD [\"./subscriber\"," + controller + "," + tunning + "," + adaptability + "," + executionType + "," + monitorInterval + "," + prefetchCount + "," + max + "," + min + "," + setPoint + "," + direction + "," + pExtra[c+t] + "]"

	fmt.Fprintf(f, "%v \n", basicDocker)
	fmt.Fprintf(f, "%v", command)
	f.Close()

	*df = append(*df, fileName)

	return
}

func createBat(list []string) {
	basicBat := "rem Self generated file at " + time.Now().String() + "\n" +
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

	listCommand := "set list="
	for i := 0; i < len(list); i++ {
		listCommand += list[i] + " "
	}
	listCommand += "\n"
	basicBat += listCommand + "\n" +
		"echo ****** BEGIN OF EXPERIMENTS *******\n" +
		"for %%x in (%list%) do (\n" +
		"echo %%x \n" +
		"   copy " + shared.DockerfilesDir + "\\" + "%%x Dockerfile\n" +
		"   docker build --tag subscriber .\n" +
		"   docker run --rm --memory=\"1g\" --cpus=\"1.0\" -v " + shared.DataDir + ":" + shared.DockerDir + " subscriber\n" +
		"   del " + shared.DockerfilesDir + "\\" + "%%x \n" +
		")\n" +
		"echo ****** END OF EXPERIMENTS *******\n"

	// create batch file
	fileName := shared.BatchFileExperiments
	f, err := os.Create(shared.BatchfilesDir + "\\" + fileName)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	fmt.Fprintf(f, "%v \n", basicBat)
	f.Close()
}
