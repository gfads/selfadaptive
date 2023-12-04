# This file has been generated automatically at 2023-12-03 20:49:11.7406636 -0300 -03 m=+0.003271001
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-execution-type=Experiment","-tunning=RootLocus","-output-file=Experiment-BasicPID-RootLocus-1","-monitor-interval=5","-direction=1.0","-hysteresis-band=200.0","-ki=2.791301426557981e-04","-controller-type=BasicPID","-alfa=1.0","-kp=-0.001395650713279","-min=1","-prefetch-count=1","-set-point=500","-dead-zone=200.0","-kd=0.001744563391599","-max=100","-beta=0.9"]