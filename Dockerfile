# syntax=docker/dockerfile:1

FROM golang:1.19

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY ./ ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
#EXPOSE 8080

ENV GOROOT=/usr/local/go/bin/

# OnOff
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=StaticGoal", "-controller-type=OnOff", "-monitor-interval=5", "-prefetch-count=1",  "-max=25", "-min=5", "-set-point=1000", "-direction=1"]
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=StaticGoal", "-controller-type=OnOffwithDeadZone", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-dead-zone=1000"]
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=OnOffwithHysteresis", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-hysteresis-band=1000"]

#************** Training - Root Locus
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=RootLocusTraining", "-controller-type=BasicP", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=0", "-direction=1","-kp=0.0", "-ki=0.0", "-kd=0.0"]
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=RootLocusTraining", "-controller-type=BasicPI", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=0", "-direction=1","-kp=0.0", "-ki=0.0", "-kd=0.0"]
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=RootLocusTraining", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=0", "-direction=1","-kp=0.0", "-ki=0.0", "-kd=0.0"]

#************** Training - Ziegler/Cohen/Amigo
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=ZieglerTraining", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=0", "-direction=1","-kp=0.0", "-ki=0.0", "-kd=0.0"]
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=ZieglerTraining", "-controller-type=BasicPI", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=0", "-direction=1","-kp=0.0", "-ki=0.0", "-kd=0.0"]

# Execution - PID - Ziegler
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00504796", "-ki=0.00252398", "-kd=0.00252398"]
CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPI", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00384517", "-ki=0.00115355", "-kd=0.00000000"]

# Execution - PID - Cohen
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00159011", "-ki=0.00150530", "-kd=0.00020288"]

# Execution - PID - AMIGO
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=0.00103062", "-ki=0.00216861", "-kd=0.00012883"]

#************** AsTAR
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=StaticGoal", "-controller-type=AsTAR", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-hysteresis-band=200"]
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=AsTAR", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-hysteresis-band=200"]

#************** Execution - PI
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPI", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=-0.00225674", "-ki=0.00217010", "-kd=0.00000000"]

#************** Execution - P
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicP", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=0.00382063", "-ki=0.00000000", "-kd=0.00000000"]

#************** Execution - PID - Root locus
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=-0.00173398", "-ki=0.00274580", "-kd=0.00063856"]
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=StaticGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-kp=-0.00173398", "-ki=0.00274580", "-kd=0.00063856"]

#************** Execution - PID - Cohen
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=0.00771191", "-ki=0.04138725", "-kd=0.00023978"]

#************** Execution - PID - AMIGO
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1",   "-kp=0.00314673", "-ki=0.02884503", "-kd=0.00012103"]

#PID Deadzone
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=DeadZonePID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=-0.00173398", "-ki=0.00274580", "-kd=0.00063856", "-dead-zone=200.0"]

#PID SmoothingDerivative
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=SmoothingDerivativePID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-kp=-0.00173398", "-ki=0.00274580", "-kd=0.00063856"]

#PID IncrementalForm
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=IncrementalFormPID", "-monitor-interval=5", "-prefetch-count=1",  "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-kp=-0.00173398", "-ki=0.00274580", "-kd=0.00063856"]

#PID ErrorSquare
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=StaticGoal", "-controller-type=ErrorSquarePID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-kp=-0.00173398", "-ki=0.00274580", "-kd=0.00063856"]

#PID GainScheduling
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=StaticGoal", "-controller-type=GainScheduling", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-kp=1.0", "-ki=1.0", "-kd=1.0", "-gain-trigger=100.0"]

#PI
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-kp=-0.00338306", "-ki=0.01070994", "-kd=0.0"]
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=StaticGoal", "-controller-type=DeadZonePID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-kp=-0.00338306", "-ki=0.01070994", "-kd=0.0", "-dead-zone=200.0"]

#CMD ["ls","-l"]
#CMD ["which","go"]