# This file has been generated automatically at 2023-10-04 16:58:24.2895584 -0300 -03 m=+0.003794501
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-controller-type=SmoothingDerivativePID","-tunning=AMIGO","-monitor-interval=5","-dead-zone=200.0","-kp=0.00053993","-output-file=Experiment-SmoothingDerivativePID-AMIGO-10","-prefetch-count=1","-kd=0.00000675","-ki=0.01136109","-execution-type=Experiment","-max=100","-set-point=500","-direction=1.0","-hysteresis-band=200.0","-beta=0.9","-min=1","-alfa=1.0"]