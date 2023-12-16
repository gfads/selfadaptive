# This file has been generated automatically at 2023-12-15 20:07:20.7790766 -0300 -03 m=+0.003668301
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-output-file=Experiment-AsTAR-t-1","-max=100","-prefetch-count=1","-direction=1.0","-controller-type=AsTAR","-beta=0.9","-alfa=1.0","-kp=1.0","-kd=1.0","-execution-type=Experiment","-min=1","-set-point=500","-dead-zone=200.0","-hysteresis-band=200.0","-tunning=None","-monitor-interval=5","-ki=1.0"]