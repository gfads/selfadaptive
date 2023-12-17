# This file has been generated automatically at 2023-12-17 15:56:30.0202282 -0300 -03 m=+0.002663401
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-execution-type=Experiment","-max=100","-dead-zone=200.0","-tunning=RootLocus","-prefetch-count=1","-direction=1.0","-kd=0.0000","-alfa=1.0","-ki=0.00080617","-output-file=Experiment-BasicPID-t-1","-min=1","-monitor-interval=5","-hysteresis-band=200.0","-controller-type=BasicPID","-set-point=500","-beta=0.9","-kp=0.00029493"]