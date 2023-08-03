# Self generated file at 2023-08-03 07:38:35.6191357 -0300 -03 m=+0.002106801
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/ 
CMD ["./subscriber","-controller-type=BasicPID","-tunning=RootLocus","-is-adaptive=true","-execution-type=Experiment","-monitor-interval=5", "-prefetch-count=1","-max=100.0","-min=1.0", "-set-point=1000", "-direction=1", "-kp=-0.000144086", "-ki=0.00248495", "-kd=0.00057789"]