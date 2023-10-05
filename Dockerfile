# This file has been generated automatically at 2023-10-05 14:43:27.7784321 -0300 -03 m=+0.004125501
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-max=100","-set-point=500","-kp=0.00026446","-ki=0.00132228","-execution-type=Experiment","-output-file=Experiment-GainScheduling-Ziegler-10","-prefetch-count=1","-direction=1.0","-tunning=Ziegler","-min=1","-monitor-interval=5","-dead-zone=200.0","-hysteresis-band=200.0","-kd=0.00001322","-controller-type=GainScheduling","-alfa=1.0","-beta=0.9"]