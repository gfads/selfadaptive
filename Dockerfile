# This file has been generated automatically at 2023-11-15 20:54:30.2157085 -0300 -03 m=+0.003383901
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-prefetch-count=1","-alfa=1.0","-ki=2.96e-05","-execution-type=Experiment","-monitor-interval=5","-max=100","-set-point=500","-direction=1.0","-dead-zone=200.0","-kd=0.0665","-controller-type=BasicPID","-tunning=RootLocus","-kp=0.0267","-output-file=Experiment-BasicPID-RootLocus-1","-hysteresis-band=200.0","-min=1","-beta=0.9"]