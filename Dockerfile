# This file has been generated automatically at 2023-09-01 15:09:53.8894862 -0300 -03 m=+0.002647901
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-min=1","-set-point=500","-direction=1.0","-hysteresis-band=200.0","-ki=0.00790004","-output-file=cohen-training-03","-max=100","-monitor-interval=5","-kp=0.00083451","-controller-type=BasicPID","-dead-zone=200.0","-beta=0.9","-alfa=1.0","-kd=0.00001065","-execution-type=Experiment","-prefetch-count=1","-tunning=Cohen"]