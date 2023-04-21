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

# Run
# OnOff
#CMD ["./subscriber","-is-adaptive=true", "-monitor-interval=5", "-prefetch-count=1", "-controller-type=OnOff", "-max=25", "-min=5", "-set-point=500", "-direction=1"]
CMD ["./subscriber","-is-adaptive=true", "-monitor-interval=5", "-prefetch-count=1", "-controller-type=OnOffwithDeadZone", "-max=10000", "-min=1", "-set-point=500", "-direction=1", "-dead-zone=1000"]
#CMD ["./subscriber","-is-adaptive=true", "-monitor-interval=5", "-prefetch-count=1", "-controller-type=OnOffwithHysteresis", "-max=25", "-min=5", "-set-point=500", "-direction=1", "-hysteresis-band=100"]

# PID
#CMD ["./subscriber","-is-adaptive=true", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=500", "-direction=1", "-kp=-0.00338306", "-ki=0.01070994", "-kd=0.00249068"]
#CMD ["./subscriber","-is-adaptive=true", "-controller-type=BasicPID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=500", "-direction=1", "-kp=-0.00338306", "-ki=0.01070994", "-kd=0.0"]
#CMD ["./subscriber","-is-adaptive=true", "-controller-type=DeadZonePID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=500", "-direction=1", "-kp=-0.00338306", "-ki=0.01070994", "-kd=0.0", "-dead-zone=200.0"]
#CMD ["./subscriber","-is-adaptive=true", "-controller-type=SmoothingDerivativePID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=500", "-direction=1", "-kp=-0.00338306", "-ki=0.01070994", "-kd=0.00249068"]
#CMD ["./subscriber","-is-adaptive=true", "-controller-type=IncrementalFormPID", "-monitor-interval=5", "-prefetch-count=1",  "-max=10000", "-min=1", "-set-point=500", "-direction=1", "-kp=-0.00338306", "-ki=0.01070994", "-kd=0.00249068"]
#CMD ["./subscriber","-is-adaptive=true", "-controller-type=DeadZonePID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=500", "-direction=1", "-kp=-0.00338306", "-ki=0.01070994", "-kd=0.00249068", "-dead-zone=200.0"]
#CMD ["./subscriber","-is-adaptive=true", "-controller-type=ErrorSquarePID", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=500", "-direction=1", "-kp=-0.00338306", "-ki=0.01070994", "-kd=0.00249068"]
#CMD ["./subscriber","-is-adaptive=true", "-controller-type=GainScheduling", "-monitor-interval=5", "-prefetch-count=1", "-max=10000", "-min=1", "-set-point=500", "-direction=1", "-kp=1.0", "-ki=1.0", "-kd=1.0", "-gain-trigger=100.0"]

#CMD ["ls","-l"]
#CMD ["which","go"]