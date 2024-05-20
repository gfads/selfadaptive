# This file has been generated automatically at 2024-05-14 17:39:41.9033849 -0300 -03 m=+0.006630401
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-kp=1.0","-ki=1.0","-kd=1.0","-tunning=None","-output-file=Experiment-Fuzzy-t-1","-monitor-interval=5","-dead-zone=200.0","-beta=0.9","-direction=1.0","-controller-type=Fuzzy","-min=1","-set-point=500","-hysteresis-band=200.0","-execution-type=Experiment","-max=100","-prefetch-count=1","-alfa=1.0"]