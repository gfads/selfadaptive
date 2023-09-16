# This file has been generated automatically at 2023-09-16 09:44:32.4748993 -0300 -03 m=+0.003663101
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-execution-type=Experiment","-direction=1.0","-beta=0.9","-kp=0.00083304","-controller-type=DeadZonePID","-max=100","-tunning=Cohen","-output-file=Experiment-DeadZonePID-Cohen-10","-monitor-interval=5","-set-point=500","-alfa=1.0","-ki=0.00788611","-min=1","-prefetch-count=1","-dead-zone=200.0","-hysteresis-band=200.0","-kd=0.00001063"]