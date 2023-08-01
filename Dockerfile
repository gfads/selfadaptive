# Self generated file at 2023-07-31 21:48:53.895597 -0300 -03 m=+0.041685201
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/ 
CMD ["./subscriber","-controller-type=SetpointWeighting","-tunning=Cohen","-is-adaptive=true", "-execution-type=DynamicGoal","-monitor-interval=5", "-prefetch-count=1","-max=100.0","-min=1.0", "-set-point=1000", "-direction=1", "-kp=0.00156457", "-ki=0.00148113", "-kd=0.00019962", "-alfa=1","-beta=0.5"]