# This file has been generated automatically at 2023-12-09 16:04:04.5751922 -0300 -03 m=+0.004242901
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-monitor-interval=5","-set-point=500","-alfa=1.0","-kd=0.0","-min=1","-max=100","-direction=1.0","-dead-zone=200.0","-hysteresis-band=200.0","-beta=0.9","-execution-type=Experiment","-output-file=Experiment-BasicPID-RootLocus-1","-ki=0.001158700597569","-controller-type=BasicPID","-tunning=RootLocus","-prefetch-count=1","-kp=-0.002896751493921"]