# This file has been generated automatically at 2023-08-13 20:45:03.5068775 -0300 -03 m=+0.003396001
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-controller-type=None","-execution-type=Static","-is-adaptive=false","-monitor-interval=5","-prefetch-count=1","-max=100","-min=1"]