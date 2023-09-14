# This file has been generated automatically at 2023-09-14 13:39:42.1831398 -0300 -03 m=+0.002641201
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-min=1","-beta=0.9","-kp=0.00053993","-kd=0.00000675","-output-file=Experiment-IncrementalFormPID-AMIGO-10","-tunning=AMIGO","-max=100","-set-point=500","-ki=0.01136109","-execution-type=Experiment","-monitor-interval=5","-dead-zone=200.0","-hysteresis-band=200.0","-controller-type=IncrementalFormPID","-direction=1.0","-alfa=1.0","-prefetch-count=1"]