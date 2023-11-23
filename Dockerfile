# This file has been generated automatically at 2023-11-23 07:54:29.2600571 -0300 -03 m=+0.004305401
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-direction=1.0","-dead-zone=200.0","-kp=1.825841182263492e-05","-ki=7.303364729053967e-06","-beta=0.9","-kd=0.0","-execution-type=Experiment","-output-file=Experiment-BasicPID-RootLocus-1","-min=1","-set-point=500","-hysteresis-band=200.0","-controller-type=BasicPID","-prefetch-count=1","-alfa=1.0","-tunning=RootLocus","-max=100","-monitor-interval=5"]