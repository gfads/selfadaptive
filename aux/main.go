package main

import (
	"fmt"
	"main.go/shared"
)

func main() {
	parameters := make(map[string]string)
	dockerBase := "FROM golang:1.19 \n" +
		"WORKDIR /app\n" +
		"COPY go.mod go.sum ./\n" +
		"RUN go mod download\n" +
		"COPY ./ ./ \n" +
		"RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go \n" +
		"ENV GOROOT=/usr/local/go/bin/"
	//tunning := []string{shared.RootLocus, shared.Ziegler, shared.Cohen, shared.Amigo}
	//controllers := []string{}
	parameters["execution-type"] = shared.DynamicGoal
	parameters["is-adaptive"] = "true"
	parameters["controller-type"] = shared.BasicPid
	parameters["monitor-interval"] = "5"
	parameters["set-point"] = "1000"
	parameters["kp"] = "0.0"
	parameters["ki"] = "0.0"
	parameters["kd"] = "0.0"
	parameters["prefetch-count"] = "1.0"
	parameters["min"] = "1.0"
	parameters["max"] = "100.0"
	parameters["dead-zone"] = "200.0"
	parameters["hysteresis-band"] = "200.0"
	parameters["direction"] = "1.0"
	parameters["gain-trigger"] = "0.0"
	parameters["alfa"] = "1.0"
	parameters["beta"] = "1.0"
	parameters["tunning"] = "RootLocus"
	command := []string{}

	command = append(command, "CMD [\"./subscriber\", \"-controller-type=OnOff\",\"-tunning=None\",\"-is-adaptive=true\", \"-execution-type=DynamicGoal\", \"-monitor-interval=5\", \"-prefetch-count=1\",  \"-max=1000\", \"-min=1\", \"-set-point=1000\", \"-direction=1\"]")

	fmt.Println(dockerBase)
	fmt.Println(command[0])
	/* #CMD ["./subscriber","-is-adaptive=true", "-execution-type=StaticGoal", "-controller-type=OnOffwithDeadZone", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-dead-zone=1000"]
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=OnOffwithHysteresis", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-hysteresis-band=1000"]

	   #************** Execution Non Adaptive
	   #CMD ["./subscriber","-is-adaptive=false", "-monitor-interval=5", "-prefetch-count=0"]

	   #************** Execution - P, PI, PID, PID Deadzone, PI with two degrees of freedom - Root locus
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=PIWithTwoDegreesOfFreedom", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=-0.00144086", "-ki=0.00248495", "-kd=0.00057789", "-beta=0.5"]
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=-0.00144086", "-ki=0.00248495", "-kd=0.00057789"]
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPI", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=-0.00211325", "-ki=0.00222392", "-kd=0.00000000"]
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=StaticGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=-0.00144086", "-ki=0.00248495", "-kd=0.00057789"]
	   #CMD["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=DeadZonePID", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=-0.00144086", "-ki=0.00248495", "-kd=0.00057789", "-dead-zone=200.0"]
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicP", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00777594", "-ki=0.00000000", "-kd=0.00000000"]

	   # Execution - P, PI, PID, PID DeadZone - Ziegler
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicP", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00596588", "-ki=0.00000000", "-kd=0.00000000"]
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=0.00496689", "-ki=0.00248344", "-kd=0.00248344"]
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPI", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=0.00406761", "-ki=0.00122028", "-kd=0.00000000"]
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=DeadZonePID", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1",   "-kp=0.00496689", "-ki=0.00248344", "-kd=0.00248344", "-dead-zone=200.0"]

	   # Execution - P, PI, PID, PID DeadZone - Cohen
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicP", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00268464", "-ki=0.00000000", "-kd=0.00000000"]
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00156457", "-ki=0.00148113", "-kd=0.00019962"]
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPI", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00414897", "-ki=0.01514702", "-kd=0.00000000"]
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=DeadZonePID", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00156457", "-ki=0.00148113", "-kd=0.00019962", "-dead-zone=200.0"]

	   # Execution - P, PI, PID, PID Deadzone - AMIGO
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00101407", "-ki=0.00213378", "-kd=0.00012676"]
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPI", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00079877", "-ki=0.00218342", "-kd=0.00000000"]
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=DeadZonePID", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00101407", "-ki=0.00213378", "-kd=0.00012676", "-dead-zone=200.0"]

	   #************** AsTAR
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=AsTAR", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-hysteresis-band=200"]

	   #************** HPA
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=HPA", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1"]

	   #************** Execution - PI (Root)
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPI", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=-0.00211325", "-ki=0.00222392", "-kd=0.00000000"]

	   #************** Execution - P (Ziegler)
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicP", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=0.00596588", "-ki=0.00000000", "-kd=0.00000000"]

	   #************** Execution - P (Cohen)
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicP", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00268464", "-ki=0.00000000", "-kd=0.00000000"]

	   #************** Execution - P (Root)
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicP", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=0.00777594", "-ki=0.00000000", "-kd=0.00000000"]

	   #************** Execution - PID - Cohen
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=0.00771191", "-ki=0.04138725", "-kd=0.00023978"]

	   #************** Execution - PID - AMIGO
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1",   "-kp=0.00314673", "-ki=0.02884503", "-kd=0.00012103"]

	   #PID SmoothingDerivative
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=SmoothingDerivativePID", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=-0.00173398", "-ki=0.00274580", "-kd=0.00063856"]

	   #PID Incremental (Root)
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=IncrementalFormPID", "-monitor-interval=5", "-prefetch-count=1",  "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=-0.00144086", "-ki=0.00248495", "-kd=0.00057789"]

	   #PID Incremental (Ziegler)
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=IncrementalFormPID", "-monitor-interval=5", "-prefetch-count=1",  "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00496689", "-ki=0.00248344", "-kd=0.00248344"]

	   #PID Incremental (Cohen)
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=IncrementalFormPID", "-monitor-interval=5", "-prefetch-count=1",  "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00156457", "-ki=0.00148113", "-kd=0.00019962", "> apague.txt"]

	   #PID Incremental (AMIGO)
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=IncrementalFormPID", "-monitor-interval=5", "-prefetch-count=1",  "-max=1000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=0.00101407", "-ki=0.00213378", "-kd=0.00012676"]

	   #PID Smoothing (Root)
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=SmoothingDerivativePID", "-monitor-interval=5", "-prefetch-count=1",  "-max=1000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=-0.00144086", "-ki=0.00248495", "-kd=0.00057789"]

	   #PID Smoothing (Ziegler)
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=SmoothingDerivativePID", "-monitor-interval=5", "-prefetch-count=1",  "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00496689", "-ki=0.00248344", "-kd=0.00248344"]

	   #PID Smoothing (Cohen) // remember it has been executed with sin(x)
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=SmoothingDerivativePID", "-monitor-interval=5", "-prefetch-count=1",  "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00156457", "-ki=0.00148113", "-kd=0.00019962"]

	   #PID Smoothing (AMIGO)
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=SmoothingDerivativePID", "-monitor-interval=5", "-prefetch-count=1",  "-max=1000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=0.00101407", "-ki=0.00213378", "-kd=0.00012676"]

	   #PID ErrorSquare
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=StaticGoal", "-controller-type=ErrorSquarePID", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=-0.00173398", "-ki=0.00274580", "-kd=0.00063856"]

	   #PID GainScheduling
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=GainScheduling", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=1.0", "-ki=1.0", "-kd=1.0", "-gain-trigger=1000.0"]

	   #PI
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=-0.00338306", "-ki=0.01070994", "-kd=0.0"]
	   #CMD ["./subscriber","-is-adaptive=true", "-execution-type=StaticGoal", "-controller-type=DeadZonePID", "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1", "-kp=-0.00338306", "-ki=0.01070994", "-kd=0.0", "-dead-zone=200.0"]

	   #CMD ["ls","-l"]
	   #CMD ["which","go"]}
	*/
}
