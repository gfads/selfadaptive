# This file has been generated automatically at 2023-09-12 16:24:05.3686581 -0300 -03 m=+0.003240501
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-direction=1.0","-beta=0.9","-kp=1.0","-ki=1.0","-controller-type=HPA","-max=100","-hysteresis-band=200.0","-min=1","-monitor-interval=5","-dead-zone=200.0","-kd=1.0","-output-file=Experiment-HPA-None-10","-tunning=None","-prefetch-count=1","-set-point=500","-alfa=1.0","-execution-type=Experiment"]