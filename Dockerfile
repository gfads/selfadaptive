# This file has been generated automatically at 2023-12-30 16:11:59.7420548 -0300 -03 m=+0.003264001
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-kd=0.00098382","-output-file=Experiment-ErrorSquarePIDFull-t-1","-monitor-interval=5","-prefetch-count=1","-alfa=1.0","-kp=-0.00174506","-ki=0.00423043","-controller-type=ErrorSquarePIDFull","-min=1","-max=100","-set-point=500","-execution-type=Experiment","-hysteresis-band=200.0","-tunning=RootLocus","-direction=1.0","-dead-zone=200.0","-beta=0.9"]