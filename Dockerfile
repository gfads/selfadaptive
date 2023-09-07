# This file has been generated automatically at 2023-09-06 21:37:59.1331565 -0300 -03 m=+0.003183901
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-dead-zone=200.0","-kp=0.00026446","-monitor-interval=5","-set-point=500","-hysteresis-band=200.0","-ki=0.00132228","-execution-type=Experiment","-controller-type=BasicPID","-max=100","-min=1","-prefetch-count=1","-beta=0.9","-alfa=1.0","-kd=0.00001322","-tunning=Ziegler","-output-file=Experiment-BasicPID-Ziegler-10","-direction=1.0"]