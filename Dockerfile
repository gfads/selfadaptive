# This file has been generated automatically at 2023-10-27 16:01:55.9494153 -0300 -03 m=+0.003351201
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-set-point=500","-alfa=1.0","-ki=0.00141098","-kd=0.00032814","-controller-type=BasicPID","-output-file=ExperimentalDesign-BasicPID-RootLocus-1","-dead-zone=200.0","-hysteresis-band=200.0","-beta=0.9","-min=1","-direction=1.0","-max=100","-monitor-interval=5","-prefetch-count=1","-kp=-0.00088937","-execution-type=ExperimentalDesign","-tunning=RootLocus"]