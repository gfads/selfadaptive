# This file has been generated automatically at 2023-08-10 12:54:09.3262562 -0300 -03 m=+0.002684501
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-controller-type=BasicPID","-tunning=RootLocus","-is-adaptive=true","-execution-type=Experiment","-monitor-interval=5","-prefetch-count=1","-max=100","-min=1","-set-point=500","-direction=1.0","-kp=-0.00144086", "-ki=0.00248495", "-kd=0.00057789"]