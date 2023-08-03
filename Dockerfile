# Self generated file at 2023-08-03 17:42:26.5381324 -0300 -03 m=+0.003483301
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/ 
CMD ["./subscriber","-controller-type=AsTAR","-tunning=None","-is-adaptive=true","-execution-type=Experiment","-monitor-interval=4","-prefetch-count=1","-max=100","-min=1","-set-point=500","-direction=1.0","-hysteresis-band=200.0"]