# This file has been generated automatically at 2023-09-17 14:31:54.5864221 -0300 -03 m=+0.003203101
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-kp=0.00053993","-controller-type=PIWithTwoDegreesOfFreedom","-set-point=500","-direction=1.0","-ki=0.01136109","-monitor-interval=5","-dead-zone=200.0","-hysteresis-band=200.0","-alfa=1.0","-kd=0.00000675","-output-file=Experiment-PIWithTwoDegreesOfFreedom-AMIGO-10","-max=100","-min=1","-prefetch-count=1","-beta=0.9","-execution-type=Experiment","-tunning=AMIGO"]