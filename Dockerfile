# This file has been generated automatically at 2023-09-03 07:21:43.1376568 -0300 -03 m=+0.002963601
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-ki=0.00218342","-direction=1.0","-kp=0.00079877","-set-point=500","-min=1","-monitor-interval=5","-prefetch-count=1","-hysteresis-band=200.0","-kd=0.0","-execution-type=Experiment","-tunning=AMIGO","-max=100","-dead-zone=200.0","-beta=0.9","-alfa=1.0","-controller-type=BasicPI","-output-file=test-1.csv"]