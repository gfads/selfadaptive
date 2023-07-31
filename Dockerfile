FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/ 
CMD ["./subscriber","-controller-type=SetpointWeighting","-tunning=AMIGO","-is-adaptive=true", "-execution-type=DynamicGoal","-monitor-interval=5", "-prefetch-count=1","-max=100.0","-min=1.0", "-set-point=1000", "-direction=1"]