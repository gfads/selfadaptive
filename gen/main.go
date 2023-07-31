package main

import (
	"fmt"
	"main.go/shared"
	"os"
)

func main() {

	list := []string{}

	/*
		// Onoff
		controllers := []string{shared.BasicOnoff, shared.DeadZoneOnoff, shared.HysteresisOnoff}
		tunnings := []string{shared.None}
		for c := 0; c < len(controllers); c++ {
			for t := 0; t < len(tunnings); t++ {
				createDockerFiles(controllers[c], tunnings[t], &list)
			}
		}
	*/
	// Astar / HPA
	controllers := []string{shared.AsTAR, shared.HPA}
	tunnings := []string{shared.None}
	for c := 0; c < len(controllers); c++ {
		for t := 0; t < len(tunnings); t++ {
			createDockerFiles(controllers[c], tunnings[t], &list)
		}
	}

	fmt.Println("Docker files created...")

	// P
	controllers = []string{shared.BasicP}
	tunnings = []string{shared.RootLocus, shared.Ziegler, shared.Cohen}
	for c := 0; c < len(controllers); c++ {
		for t := 0; t < len(tunnings); t++ {
			createDockerFiles(controllers[c], tunnings[t], &list)
		}
	}

	// PI
	controllers = []string{shared.BasicPi}
	tunnings = []string{shared.RootLocus, shared.Ziegler, shared.Cohen, shared.Amigo}
	for c := 0; c < len(controllers); c++ {
		for t := 0; t < len(tunnings); t++ {
			createDockerFiles(controllers[c], tunnings[t], &list)
		}
	}
	// PID
	controllers = []string{shared.BasicPid, shared.DeadZonePid, shared.IncrementalFormPid, shared.SmoothingPid, shared.GainScheduling, shared.SetpointWeighting, shared.SetpointWeighting}
	tunnings = []string{shared.RootLocus, shared.Ziegler, shared.Cohen, shared.Amigo}
	for c := 0; c < len(controllers); c++ {
		for t := 0; t < len(tunnings); t++ {
			createDockerFiles(controllers[c], tunnings[t], &list)
		}
	}

	// create execute
	createBat(list)

	fmt.Println("Bat file created...")

	/*
		wg := sync.WaitGroup{}
		// execute exec

			wg.Add(1)
			cmd := exec.Command("C:\\Users\\user\\go\\selfadaptive\\Apague-execute.bat")
			go func() {
				defer wg.Done()
				err := cmd.Run()
				if err != nil {
					log.Fatal(err)
				}
				cmd.Wait()
			}()
			wg.Wait()
			fmt.Println("Bat file executed...")
	*/

}

func createDockerFiles(c, t string, list *[]string) {

	// extra parameters
	pExtra := make(map[string]string)
	pExtra[shared.BasicOnoff+shared.None] = ""
	pExtra[shared.DeadZoneOnoff+shared.None] = ", \"-dead-zone=200\""
	pExtra[shared.HysteresisOnoff+shared.None] = ", \"-hysteresis-band=200\""
	pExtra[shared.AsTAR+shared.None] = ", \"-hysteresis-band=200\""
	pExtra[shared.BasicP+shared.RootLocus] = ", \"-kp=0.00777594\", \"-ki=0.00000000\", \"-kd=0.00000000\""
	pExtra[shared.BasicP+shared.Ziegler] = ", \"-kp=0.00596588\", \"-ki=0.00000000\", \"-kd=0.00000000\""
	pExtra[shared.BasicP+shared.Cohen] = ", \"-kp=0.00268464\", \"-ki=0.00000000\", \"-kd=0.00000000\""
	pExtra[shared.BasicPi+shared.RootLocus] = ", \"-kp=-0.00211325\", \"-ki=0.00222392\", \"-kd=0.00000000\""
	pExtra[shared.BasicPi+shared.Ziegler] = ", \"-kp=0.00406761\", \"-ki=0.00122028\", \"-kd=0.00000000\""
	pExtra[shared.BasicPi+shared.Cohen] = ", \"-kp=0.00414897\", \"-ki=0.01514702\", \"-kd=0.00000000\""
	pExtra[shared.BasicPi+shared.Amigo] = ", \"-kp=0.00079877\", \"-ki=0.00218342\", \"-kd=0.00000000\""
	pExtra[shared.BasicPid+shared.RootLocus] = ", \"-kp=-0.00144086\", \"-ki=0.00248495\", \"-kd=0.00057789\""
	pExtra[shared.BasicPid+shared.Ziegler] = ", \"-kp=0.00496689\", \"-ki=0.00248344\", \"-kd=0.00248344\""
	pExtra[shared.BasicPid+shared.Cohen] = ", \"-kp=0.00156457\", \"-ki=0.00148113\", \"-kd=0.00019962\""
	pExtra[shared.BasicPid+shared.Amigo] = ", \"-kp=0.00101407\", \"-ki=0.00213378\", \"-kd=0.00012676\""
	pExtra[shared.DeadZonePid+shared.RootLocus] = ", \"-kp=-0.00144086\", \"-ki=0.00248495\", \"-kd=0.00057789\", \"-dead-zone=200\""
	pExtra[shared.DeadZonePid+shared.Ziegler] = ", \"-kp=0.00496689\", \"-ki=0.00248344\", \"-kd=0.00248344\", \"-dead-zone=200\""
	pExtra[shared.DeadZonePid+shared.Cohen] = ", \"-kp=0.00156457\", \"-ki=0.00148113\", \"-kd=0.00019962\", \"-dead-zone=200\""
	pExtra[shared.DeadZonePid+shared.Amigo] = ", \"-kp=0.00101407\", \"-ki=0.00213378\", \"-kd=0.00012676\", \"-dead-zone=200\""
	pExtra[shared.IncrementalFormPid+shared.RootLocus] = ", \"-kp=-0.00144086\", \"-ki=0.00248495\", \"-kd=0.00057789\""
	pExtra[shared.IncrementalFormPid+shared.Ziegler] = ", \"-kp=0.00496689\", \"-ki=0.00248344\", \"-kd=0.00248344\""
	pExtra[shared.IncrementalFormPid+shared.Cohen] = ", \"-kp=0.00156457\", \"-ki=0.00148113\", \"-kd=0.00019962\""
	pExtra[shared.IncrementalFormPid+shared.Amigo] = ", \"-kp=0.00101407\", \"-ki=0.00213378\", \"-kd=0.00012676\""
	pExtra[shared.SmoothingPid+shared.RootLocus] = ", \"-kp=-0.00144086\", \"-ki=0.00248495\", \"-kd=0.00057789\""
	pExtra[shared.SmoothingPid+shared.Ziegler] = ", \"-kp=0.00496689\", \"-ki=0.00248344\", \"-kd=0.00248344\""
	pExtra[shared.SmoothingPid+shared.Cohen] = ", \"-kp=0.00156457\", \"-ki=0.00148113\", \"-kd=0.00019962\""
	pExtra[shared.SmoothingPid+shared.Amigo] = ", \"-kp=0.00101407\", \"-ki=0.00213378\", \"-kd=0.00012676\""
	pExtra[shared.PIwithTwoDegreesOfFreedom+shared.RootLocus] = ", \"-kp=-0.00211325\", \"-ki=0.00222392\", \"-kd=0.00000000\",\"-beta=0.5\""
	pExtra[shared.PIwithTwoDegreesOfFreedom+shared.Ziegler] = ", \"-kp=0.00406761\", \"-ki=0.00122028\", \"-kd=0.00000000\", \"-beta=0.5\""
	pExtra[shared.PIwithTwoDegreesOfFreedom+shared.Cohen] = ", \"-kp=0.00414897\", \"-ki=0.01514702\", \"-kd=0.00000000\", \"-beta=0.5\""
	pExtra[shared.PIwithTwoDegreesOfFreedom+shared.Amigo] = ", \"-kp=0.00079877\", \"-ki=0.00218342\", \"-kd=0.00000000\", \"-beta=0.5\""
	pExtra[shared.SetpointWeighting+shared.RootLocus] = ", \"-kp=-0.00144086\", \"-ki=0.00248495\", \"-kd=0.00057789\", \"-alfa=1\",\"-beta=0.5\""
	pExtra[shared.SetpointWeighting+shared.Ziegler] = ", \"-kp=0.00496689\", \"-ki=0.00248344\", \"-kd=0.00248344\", \"-alfa=1\",\"-beta=0.5\""
	pExtra[shared.SetpointWeighting+shared.Cohen] = ", \"-kp=0.00156457\", \"-ki=0.00148113\", \"-kd=0.00019962\", \"-alfa=1\",\"-beta=0.5\""
	pExtra[shared.SetpointWeighting+shared.Amigo] = ", \"-kp=0.00101407\", \"-ki=0.00213378\", \"-kd=0.00012676\", \"-alfa=1\",\"-beta=0.5\""

	// docker files folder
	dir := "C:\\Users\\user\\go\\selfadaptive"
	basicDocker := "FROM golang:1.19\n" +
		"WORKDIR /app\n" +
		"COPY go.mod go.sum ./\n" +
		"RUN go mod download \n" +
		"COPY ./ ./ \n" +
		"RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go\n" +
		"ENV GOROOT=/usr/local/go/bin/"

	min := "\"-min=1.0\""
	max := "\"-max=100.0\""
	monitorInterval := "\"-monitor-interval=5\""

	controller := "\"-controller-type=" + c + "\""
	tunning := "\"-tunning=" + t + "\""
	fileName := "Dockerfile-" + c + "-" + t
	f, err := os.Create(dir + "\\" + fileName)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	*list = append(*list, fileName)

	command := "CMD [\"./subscriber\"," + controller + "," + tunning + ",\"-is-adaptive=true\", \"-execution-type=DynamicGoal\"," + monitorInterval + ", \"-prefetch-count=1\"," + max + "," + min + ", \"-set-point=1000\", \"-direction=1\"" + pExtra[c+t] + "]"

	fmt.Fprintf(f, "%v \n", basicDocker)
	fmt.Fprintf(f, "%v", command)
	f.Close()
	return
}

func createBat(list []string) {
	basicBat := "@echo off \n" +
		"docker stop some-rabbit \n" +
		"docker rm some-rabbit\n" +
		"docker run -d --memory=\"6g\" --cpus=\"5.0\" --name some-rabbit -p 5672:5672 rabbitmq\n" +
		"timeout /t 10\n" +
		"@echo Removing previous containers\n" +
		"docker stop publisher\n" +
		"docker rm publisher\n" +
		"docker stop subscriber\n" +
		"docker rm subscriber\n"

	listCommand := "set list="
	for i := 0; i < len(list); i++ {
		listCommand += list[i] + " "
	}
	listCommand += "\n"
	basicBat += listCommand + "echo ****** BEGIN OF EXPERIMENTS *******\n" +
		"for %%x in (%list%) do (\n" +
		"	echo %%x\n" +
		"   copy %%x Dockerfile\n" +
		"   docker build --tag subscriber .\n" +
		"	docker run --memory=\"1g\" --cpus=\"1.0\" -v C:\\Users\\user\\go\\selfadaptive\\rabbitmq\\data:/app/data subscriber\n" +
		")\n" +
		"echo ****** END OF EXPERIMENTS *******\n"

	// docker files folder
	dir := "C:\\Users\\user\\go\\selfadaptive"
	fileName := "Apague-execute.bat"
	f, err := os.Create(dir + "\\" + fileName)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	fmt.Fprintf(f, "%v \n", basicBat)
	f.Close()
}
