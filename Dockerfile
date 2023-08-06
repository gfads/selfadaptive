# This file has been generated automatically at 2023-08-06 09:08:32.7880222 -0300 -03 m=+0.019538201
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-controller-type=PIWithTwoDegreesOfFreedom","-tunning=AMIGO","-is-adaptive=true","-execution-type=Experiment","-monitor-interval=5","-prefetch-count=1","-max=100","-min=1","-set-point=500","-direction=1.0",,"-beta=+0.5"]