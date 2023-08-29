# This file has been generated automatically at 2023-08-29 14:01:28.3373713 -0300 -03 m=+0.003151001
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-alfa=1.0","-kp=0.00031917","-kd=0.00000478","-set-point=500","-dead-zone=200.0","-beta=0.9","-execution-type=Experiment","-controller-type=BasicPID","-output-file=test.csv","-max=100","-monitor-interval=5","-min=1","-prefetch-count=1","-hysteresis-band=200.0","-tunning=Ziegler","-direction=1.0","-ki=0.00047798"]