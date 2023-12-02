# This file has been generated automatically at 2023-12-02 09:56:15.0450407 -0300 -03 m=+0.003544601
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-min=1","-kp=-0.001395650713279","-ki=1.165570481688410e-05","-tunning=RootLocus","-direction=1.0","-dead-zone=200.0","-alfa=1.0","-monitor-interval=5","-hysteresis-band=200.0","-execution-type=Experiment","-output-file=Experiment-BasicPID-RootLocus-1","-max=100","-prefetch-count=1","-set-point=500","-beta=0.9","-kd=7.284815510552564e-05","-controller-type=BasicPID"]