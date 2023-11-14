# This file has been generated automatically at 2023-11-13 21:18:42.5469263 -0300 -03 m=+0.004402101
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-dead-zone=200.0","-hysteresis-band=200.0","-ki=0.00141098","-controller-type=BasicPID","-prefetch-count=1","-min=1","-monitor-interval=5","-direction=1.0","-beta=0.9","-kd=0.00032814","-execution-type=ExperimentalDesign","-output-file=ExperimentalDesign-BasicPID-RootLocus-1","-alfa=1.0","-kp=-0.00088937","-set-point=500","-tunning=RootLocus","-max=100"]