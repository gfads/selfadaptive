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
	createBat(dockerFile)
}

func configureDockerFile(et, ct, t, f string) map[string]string {

	// header - common to all executions
	h := "# This file has been generated automatically at " + time.Now().String() + "\n" +
		"FROM golang:1.19\n" +
		"WORKDIR /app\n" +
		"COPY go.mod go.sum ./\n" +
		"RUN go mod download \n" +
		"COPY ./ ./ \n" +
		"RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go\n" +
		"ENV GOROOT=/usr/local/go/bin/\n"

	// command general parameters
	params := map[string]string{}
	params["min"] = shared.MinPC
	params["max"] = shared.MaxPC
	params["monitor-interval"] = shared.MonitorInterval
	params["prefetch-count"] = shared.InitialPC
	params["execution-type"] = et
	params["set-point"] = shared.SetPoint
	params["direction"] = shared.Direction
	params["tunning"] = t
	params["controller-type"] = ct
	params["output-file"] = f

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

func experimentExecution(c []string, t []string, d string) {
	dockerFileNames := []string{}

	// configure dockerfile dir
	shared.ConfigureDockerfileDir(d)

	// create dockerfiles for experiments
	for i := 0; i < len(c); i++ {
		for j := 0; j < len(t); j++ {
			pExtra := loadExtraParameters(c, t)
			dockerFileNames = append(dockerFileNames, createDockerFileExperiment(c[i], t[j], pExtra))
		}
	}

	// create execute
	//createBat(dockerFileNames)
}

func staticExecution(d string) {
	dockerFileNames := []string{}

	// configure dockerfile dir
	shared.ConfigureDockerfileDir(d)

	// create docker file training
	dockerFileNames = append(dockerFileNames, createDockerFileStatic(shared.DockerFileStatic))

	// create execute
	//createBat(dockerFileNames)
}

func rootLocusExecution(d string) {
	dockerFileNames := []string{}

	// configure dockerfile dir
	shared.ConfigureDockerfileDir(d)

	// create docker file training
	dockerFileNames = append(dockerFileNames, createDockerFileRoot(shared.DockerFileRoot))

	// create execute
	//createBat(dockerFileNames)
}

func zieglerExecution(d string) {
	dockerFileNames := []string{}

	// configure dockerfile dir
	shared.ConfigureDockerfileDir(d)

	// create docker file training
	dockerFileNames = append(dockerFileNames, createDockerFileZiegler(shared.DockerFileZiegler))

	// create execute
	//createBat(dockerFileNames)
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
			} else if c[i] == shared.PIwithTwoDegreesOfFreedom {
				keyGain = shared.BasicPi + t[j]
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

	r[shared.PIwithTwoDegreesOfFreedom+shared.RootLocus] += ",\"-beta=" + shared.Beta + "\""
	r[shared.PIwithTwoDegreesOfFreedom+shared.Ziegler] += ",\"-beta=" + shared.Beta + "\""
	r[shared.PIwithTwoDegreesOfFreedom+shared.Cohen] += ",\"-beta=" + shared.Beta + "\""
	r[shared.PIwithTwoDegreesOfFreedom+shared.Amigo] += ",\"-beta=" + shared.Beta + "\""
	r[shared.SetpointWeighting+shared.RootLocus] += ",\"-alfa=" + shared.Alfa + "\",\"-beta=" + shared.Beta + "\""
	r[shared.SetpointWeighting+shared.Ziegler] += ",\"-alfa=" + shared.Alfa + "\",\"-beta=" + shared.Beta + "\""
	r[shared.SetpointWeighting+shared.Cohen] += ",\"-alfa=" + shared.Alfa + "\",\"-beta=" + shared.Beta + "\""
	r[shared.SetpointWeighting+shared.Amigo] += ",\"-alfa=" + shared.Alfa + "\",\"-beta=" + shared.Beta + "\""

	return r
}

func createDockerFileExperiment(c, t string, pExtra map[string]string) string {
	// set tunning of non pid controllers to None
	if c == shared.AsTAR || c == shared.HPA || c == shared.BasicOnoff ||
		c == shared.DeadZoneOnoff || c == shared.HysteresisOnoff {
		t = shared.None
	}

	// configure general "command" parameters
	min := "\"-min=" + shared.MinPC + "\""
	max := "\"-max=" + shared.MaxPC + "\""
	monitorInterval := "\"-monitor-interval=" + shared.MonitorInterval + "\""
	executionType := "\"-execution-type=" + shared.Experiment + "\""
	prefetchCount := "\"-prefetch-count=" + shared.InitialPC + "\""
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
	dockerFileContent := "# This file has been generated automatically at " + time.Now().String() + "\n" +
		"FROM golang:1.19\n" +
		"WORKDIR /app\n" +
		"COPY go.mod go.sum ./\n" +
		"RUN go mod download \n" +
		"COPY ./ ./ \n" +
		"RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go\n" +
		"ENV GOROOT=/usr/local/go/bin/\n"

	// define Docker file CMD
	cmd := ""
	if pExtra[c+t] == "" {
		//cmd = "CMD [\"./subscriber\"," + controller + "," + tunning + "," + adaptability + "," + executionType + "," + monitorInterval + "," + prefetchCount + "," + max + "," + min + "," + setPoint + "," + direction + "]"
		cmd = "CMD [\"./subscriber\"," + controller + "," + tunning + "," + "," + executionType + "," + monitorInterval + "," + prefetchCount + "," + max + "," + min + "," + setPoint + "," + direction + "]"
	} else {
		//cmd = "CMD [\"./subscriber\"," + controller + "," + tunning + "," + adaptability + "," + executionType + "," + monitorInterval + "," + prefetchCount + "," + max + "," + min + "," + setPoint + "," + direction + "," + pExtra[c+t] + "]"
		cmd = "CMD [\"./subscriber\"," + controller + "," + tunning + "," + executionType + "," + monitorInterval + "," + prefetchCount + "," + max + "," + min + "," + setPoint + "," + direction + "," + pExtra[c+t] + "]"
	}

	// include command in docker file
	dockerFileContent += cmd

	// write in docker file
	fmt.Fprintf(df, dockerFileContent)

	return fileName
}

func createDockerFileStatic(fileName string) string {

	// configure general "command" parameters
	controllerType := "\"-controller-type=" + shared.None + "\""
	min := "\"-min=" + shared.MinPC + "\""
	max := "\"-max=" + shared.MaxPC + "\""
	monitorInterval := "\"-monitor-interval=" + shared.MonitorInterval + "\""
	executionType := "\"-execution-type=" + shared.InputStep + "\""
	prefetchCount := "\"-prefetch-count=" + shared.InitialPC + "\""

	// create docker file
	df, err := os.Create(shared.DockerfilesDir + "\\" + fileName)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	defer df.Close()

	// docker file content
	dockerFileContent := "# This file has been generated automatically at " + time.Now().String() + "\n" +
		"FROM golang:1.19\n" +
		"WORKDIR /app\n" +
		"COPY go.mod go.sum ./\n" +
		"RUN go mod download \n" +
		"COPY ./ ./ \n" +
		"RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go\n" +
		"ENV GOROOT=/usr/local/go/bin/\n"

	//cmd := "CMD [\"./subscriber\"," + controllerType + "," + executionType + "," + adaptability + "," + monitorInterval + "," + prefetchCount + "," + max + "," + min + "]"
	cmd := "CMD [\"./subscriber\"," + controllerType + "," + executionType + "," + monitorInterval + "," + prefetchCount + "," + max + "," + min + "]"
	// include command in docker file
	dockerFileContent += cmd

	// write in docker file
	fmt.Fprintf(df, dockerFileContent)

	return fileName
}

func createDockerFileRoot(fileName string) string {

	// configure general "command" parameters
	controllerType := "\"-controller-type=" + shared.BasicPid + "\""
	min := "\"-min=" + shared.MinPC + "\""
	max := "\"-max=" + shared.MaxPC + "\""
	monitorInterval := "\"-monitor-interval=" + shared.MonitorInterval + "\""
	executionType := "\"-execution-type=" + shared.RootTraining + "\""
	prefetchCount := "\"-prefetch-count=" + shared.InitialPC + "\""

	// create docker file
	df, err := os.Create(shared.DockerfilesDir + "\\" + fileName)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	defer df.Close()

	// docker file content
	dockerFileContent := "# This file has been generated automatically at " + time.Now().String() + "\n" +
		"FROM golang:1.19\n" +
		"WORKDIR /app\n" +
		"COPY go.mod go.sum ./\n" +
		"RUN go mod download \n" +
		"COPY ./ ./ \n" +
		"RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go\n" +
		"ENV GOROOT=/usr/local/go/bin/\n"

	//cmd := "CMD [\"./subscriber\"," + controllerType + "," + executionType + "," + adaptability + "," + monitorInterval + "," + prefetchCount + "," + max + "," + min + "]"
	cmd := "CMD [\"./subscriber\"," + controllerType + "," + executionType + "," + monitorInterval + "," + prefetchCount + "," + max + "," + min + "]"

	// include command in docker file
	dockerFileContent += cmd

	// write in docker file
	fmt.Fprintf(df, dockerFileContent)

	return fileName
}

func createDockerFileZiegler(fileName string) string {

	// configure general "command" parameters
	controllerType := "\"-controller-type=" + shared.BasicPid + "\""
	min := "\"-min=" + shared.MinPC + "\""
	max := "\"-max=" + shared.MaxPC + "\""
	monitorInterval := "\"-monitor-interval=" + shared.MonitorInterval + "\""
	executionType := "\"-execution-type=" + shared.ZieglerTraining + "\""
	prefetchCount := "\"-prefetch-count=" + shared.InitialPC + "\""

	// create docker file
	df, err := os.Create(shared.DockerfilesDir + "\\" + fileName)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	defer df.Close()

	// docker file content
	dockerFileContent := "# This file has been generated automatically at " + time.Now().String() + "\n" +
		"FROM golang:1.19\n" +
		"WORKDIR /app\n" +
		"COPY go.mod go.sum ./\n" +
		"RUN go mod download \n" +
		"COPY ./ ./ \n" +
		"RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go\n" +
		"ENV GOROOT=/usr/local/go/bin/\n"

	cmd := "CMD [\"./subscriber\"," + controllerType + "," + executionType + "," + monitorInterval + "," + prefetchCount + "," + max + "," + min + "]"

	// include command in docker file
	dockerFileContent += cmd

	// write in docker file
	fmt.Fprintf(df, dockerFileContent)

	return fileName
}

func createBat(list map[string]string) {
	// define list of docker files include in bat file
	listCommand := "set list="
	for k, _ := range list {
		listCommand += k + " "
	}
	listCommand += "\n"

	batFileContent := "rem This file has been generated automatically at " + time.Now().String() + "\n" +
		"@echo off \n" +
		"docker stop some-rabbit \n" +
		"docker rm some-rabbit\n" +
		"docker run -d --memory=\"6g\" --cpus=\"5.0\" --name some-rabbit -p 5672:5672 rabbitmq\n" +
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
			"   echo y | docker volume prune \n" +
			"   echo y | docker image prune \n" +
			")\n" +
			"docker stop some-rabbit \n" +
			"docker rm some-rabbit\n"

	// create batch file
	fileName := shared.BatchFileExperiments
	bf, err := os.Create(shared.BatchfilesDir + "\\" + fileName)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	defer bf.Close()

	fmt.Fprintf(bf, "%v", batFileContent)
}
