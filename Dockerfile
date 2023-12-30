# This file has been generated automatically at 2023-12-30 12:06:07.5879665 -0300 -03 m=+0.004779301
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-min=1","-monitor-interval=5","-prefetch-count=1","-dead-zone=200.0","-hysteresis-band=200.0","-controller-type=BasicPID","-max=100","-alfa=1.0","-execution-type=Experiment","-set-point=500","-beta=0.9","-kp=-0.001174999167787","-ki=4.699996671148544e-04","-tunning=RootLocus","-direction=1.0","-kd=0.0","-output-file=Experiment-BasicPID-t-1"]