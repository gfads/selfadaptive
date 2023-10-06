# This file has been generated automatically at 2023-10-06 10:35:40.191066 -0300 -03 m=+0.003504301
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-kp=0.00083304","-kd=0.00001063","-alfa=1.0","-execution-type=Experiment","-output-file=Experiment-ErrorSquarePIDFull-Cohen-10","-direction=1.0","-hysteresis-band=200.0","-ki=0.00788611","-max=100","-prefetch-count=1","-set-point=500","-beta=0.9","-dead-zone=200.0","-controller-type=ErrorSquarePIDFull","-tunning=Cohen","-min=1","-monitor-interval=5"]