# This file has been generated automatically at 2023-10-08 21:22:23.926093 -0300 -03 m=+0.004125401
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-execution-type=Experiment","-hysteresis-band=200.0","-beta=0.9","-kd=0.0","-controller-type=BasicPI","-tunning=AMIGO","-output-file=Experiment-BasicPI-AMIGO-10","-max=100","-direction=1.0","-kp=0.00037871","-ki=0.01035190","-min=1","-monitor-interval=5","-prefetch-count=1","-set-point=500","-dead-zone=200.0","-alfa=1.0"]