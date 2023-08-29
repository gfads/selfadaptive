# This file has been generated automatically at 2023-08-29 12:48:57.9728329 -0300 -03 m=+0.002685401
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-set-point=500","-tunning=None","-output-file=test.csv","-direction=1.0","-controller-type=BasicPID","-min=1","-max=100","-monitor-interval=5","-prefetch-count=1","-execution-type=Experiment"]