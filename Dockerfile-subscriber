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
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=ZieglerTraining", "-controller-type=BasicP", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=0", "-direction=1","-kp=0.0", "-ki=0.0", "-kd=0.0"]
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=ZieglerTraining", "-controller-type=BasicPI", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=0", "-direction=1","-kp=0.0", "-ki=0.0", "-kd=0.0"]
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=ZieglerTraining", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=0", "-direction=1","-kp=0.0", "-ki=0.0", "-kd=0.0"]

# Execution - PID (Deadzone)
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=DeadZonePID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=-0.00144086", "-ki=0.00248495", "-kd=0.00057789", "-dead-zone=200.0"]

# Execution - P, PI, PID, PID DeadZone - Ziegler
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=0.00496689", "-ki=0.00248344", "-kd=0.00248344"]
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPI", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00384517", "-ki=0.00115355", "-kd=0.00000000"]
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=DeadZonePID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1",   "-kp=0.00496689", "-ki=0.00248344", "-kd=0.00248344", "-dead-zone=200.0"]


# Execution - P, PI, PID, PID DeadZone - Cohen
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00156457", "-ki=0.00148113", "-kd=0.00019962"]
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPI", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00414897", "-ki=0.01514702", "-kd=0.00000000"]
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=DeadZonePID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00156457", "-ki=0.00148113", "-kd=0.00019962", "-dead-zone=200.0"]

# Execution - P, PI, PID, PID Deadzone - AMIGO
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00101407", "-ki=0.00213378", "-kd=0.00012676"]
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPI", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00069889", "-ki=0.00191040", "-kd=0.00000000"]
CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=DeadZonePID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-kp=0.00101407", "-ki=0.00213378", "-kd=0.00012676", "-dead-zone=200.0"]

#************** AsTAR
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=AsTAR", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-hysteresis-band=200"]

#************** HPA
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=HPA", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1"]

#************** Execution - PI (Root)
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPI", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=-0.00211325", "-ki=0.00222392", "-kd=0.00000000"]

#************** Execution - P
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicP", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=0.00382063", "-ki=0.00000000", "-kd=0.00000000"]

#************** Execution - PID - Root locus
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=-0.00144086", "-ki=0.00248495", "-kd=0.00057789"]
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=StaticGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1", "-kp=-0.00144086", "-ki=0.00248495", "-kd=0.00057789"]

#************** Execution - PID - Cohen
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=0.00771191", "-ki=0.04138725", "-kd=0.00023978"]

#************** Execution - PID - AMIGO
#CMD ["./subscriber","-is-adaptive=true", "-execution-type=DynamicGoal", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=1000", "-direction=1",   "-kp=0.00314673", "-ki=0.02884503", "-kd=0.00012103"]

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