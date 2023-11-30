# This file has been generated automatically at 2023-11-24 17:06:03.3335908 -0300 -03 m=+0.003133401
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-kd=0.00032814","-execution-type=Experiment","-output-file=Experiment-BasicPID-RootLocus-1","-set-point=500","-ki=0.00141098","-dead-zone=200.0","-hysteresis-band=200.0","-kp=-0.00088937","-tunning=RootLocus","-min=1","-prefetch-count=1","-direction=1.0","-controller-type=BasicPID","-max=100","-alfa=1.0","-monitor-interval=5","-beta=0.9"]