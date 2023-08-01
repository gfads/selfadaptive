# Self generated file at 2023-08-01 11:31:35.5615468 -0300 -03 m=+0.013804601
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go
ENV GOROOT=/usr/local/go/bin/ 
CMD ["./subscriber","-controller-type=PIWithTwoDegreesOfFreedom","-tunning=AMIGO","-is-adaptive=true", "-execution-type=DynamicGoal","-monitor-interval=5", "-prefetch-count=1","-max=100.0","-min=1.0", "-set-point=1000", "-direction=1", "-kp=0.00079877", "-ki=0.00218342", "-kd=0.00000000", "-beta=0.5"]