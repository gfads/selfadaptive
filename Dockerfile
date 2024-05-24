# This file has been generated automatically at 2024-05-23 14:12:37.9038946 -0300 -03 m=+0.003157201
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-tunning=None","-min=1","-dead-zone=200.0","-ki=1.0","-kd=1.0","-execution-type=Experiment","-controller-type=AsTAR","-prefetch-count=1","-hysteresis-band=200.0","-kp=1.0","-max=100","-monitor-interval=5","-direction=1.0","-beta=0.9","-alfa=1.0","-output-file=Experiment-AsTAR-t-2","-set-point=500"]