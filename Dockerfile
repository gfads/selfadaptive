# This file has been generated automatically at 2023-09-19 10:44:33.2527174 -0300 -03 m=+0.002666901
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-min=1","-direction=1.0","-beta=0.9","-kd=0.00000675","-controller-type=SetpointWeighting","-tunning=AMIGO","-set-point=500","-hysteresis-band=200.0","-alfa=1.0","-ki=0.01136109","-max=100","-prefetch-count=1","-kp=0.00053993","-execution-type=Experiment","-output-file=Experiment-SetpointWeighting-AMIGO-10","-monitor-interval=5","-dead-zone=200.0"]