# This file has been generated automatically at 2023-10-05 11:15:25.9207396 -0300 -03 m=+0.003162001
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-max=100","-set-point=500","-kp=0.00083304","-controller-type=DeadZonePID","-tunning=Cohen","-monitor-interval=5","-prefetch-count=1","-dead-zone=200.0","-beta=0.9","-alfa=1.0","-ki=0.00788611","-execution-type=Experiment","-min=1","-direction=1.0","-hysteresis-band=200.0","-kd=0.00001063","-output-file=Experiment-DeadZonePID-Cohen-10"]