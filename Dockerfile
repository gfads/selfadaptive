FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber", "-controller-type=BasicPID", "-tunning=AMIGO", "-is-adaptive=true", "-execution-type=DynamicGoal",  "-monitor-interval=5", "-prefetch-count=1", "-max=1000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=0.00101407", "-ki=0.00213378", "-kd=0.00012676"]
