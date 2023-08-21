# This file has been generated automatically at 2023-08-21 15:57:13.6878574 -0300 -03 m=+0.003876001
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-controller-type=BasicPID","-execution-type=ZieglerTraining","-monitor-interval=5","-prefetch-count=1","-max=100","-min=1"]