# This file has been generated automatically at 2023-12-15 16:11:31.3435581 -0300 -03 m=+0.003219001
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-controller-type=BasicPID","-output-file=Experiment-BasicPID-t-1","-min=1","-beta=0.9","-alfa=1.0","-execution-type=Experiment","-set-point=500","-hysteresis-band=200.0","-kp=-0.001877444539468","-kd=0.002346805674335","-tunning=RootLocus","-max=100","-dead-zone=200.0","-ki=3.754889078936145e-04","-monitor-interval=5","-prefetch-count=1","-direction=1.0"]