# This file has been generated automatically at 2023-08-05 20:10:02.4986883 -0300 -03 m=+0.021044901
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-controller-type=HPA","-tunning=None","-is-adaptive=true","-execution-type=Experiment","-monitor-interval=5","-prefetch-count=1","-max=100","-min=1","-set-point=500","-direction=1.0"]