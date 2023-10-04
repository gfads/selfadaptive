# This file has been generated automatically at 2023-10-04 14:34:23.8370651 -0300 -03 m=+0.003211601
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-beta=0.9","-tunning=RootLocus","-output-file=Experiment-SmoothingDerivativePID-RootLocus-10","-prefetch-count=1","-alfa=1.0","-kp=-0.00088937","-ki=0.00141098","-monitor-interval=5","-set-point=500","-direction=1.0","-dead-zone=200.0","-hysteresis-band=200.0","-execution-type=Experiment","-controller-type=SmoothingDerivativePID","-max=100","-kd=0.00032814","-min=1"]