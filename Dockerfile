# This file has been generated automatically at 2023-12-12 07:49:25.6371813 -0300 -03 m=+0.011043601
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-controller-type=BasicPID","-set-point=500","-alfa=1.0","-tunning=RootLocus","-direction=1.0","-kp=-9.763052804675233e-05","-kd=0.0","-beta=0.9","-ki=3.905221121870094e-05","-execution-type=Experiment","-monitor-interval=5","-prefetch-count=1","-dead-zone=200.0","-output-file=Experiment-BasicPID-RootLocus-1","-min=1","-max=100","-hysteresis-band=200.0"]