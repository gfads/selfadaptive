# This file has been generated automatically at 2023-10-06 14:58:20.1174163 -0300 -03 m=+0.003361901
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-kp=-0.00088937","-tunning=RootLocus","-monitor-interval=5","-ki=0.00141098","-max=100","-prefetch-count=1","-output-file=Experiment-SmoothingDerivativePID-RootLocus-10","-min=1","-set-point=500","-dead-zone=200.0","-hysteresis-band=200.0","-kd=0.00032814","-execution-type=Experiment","-controller-type=SmoothingDerivativePID","-alfa=1.0","-direction=1.0","-beta=0.9"]