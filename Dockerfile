# This file has been generated automatically at 2023-11-23 14:55:23.3618779 -0300 -03 m=+0.002809301
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-controller-type=BasicPID","-output-file=Experiment-BasicPID-RootLocus-1","-kp=-0.00060587","-min=1","-direction=1.0","-dead-zone=200.0","-prefetch-count=3","-hysteresis-band=200.0","-ki=0.00012117","-monitor-interval=5","-set-point=500","-beta=0.9","-alfa=1.0","-kd=0.00075734","-execution-type=Experiment","-tunning=RootLocus","-max=100"]