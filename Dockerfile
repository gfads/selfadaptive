# This file has been generated automatically at 2023-11-14 11:04:07.0506177 -0300 -03 m=+0.003290501
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-execution-type=Experiment","-controller-type=HPA","-output-file=Experiment-HPA-None-1","-min=1","-monitor-interval=5","-prefetch-count=1","-direction=1.0","-dead-zone=200.0","-kd=1.0","-set-point=500","-hysteresis-band=200.0","-max=100","-beta=0.9","-kp=1.0","-ki=1.0","-tunning=None","-alfa=1.0"]