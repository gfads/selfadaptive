# syntax=docker/dockerfile:1

FROM golang:1.19

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY ./ ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o ./subscriber ./rabbitmq/subscriber/main.go

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
#EXPOSE 8080

ENV GOROOT=/usr/local/go/bin/
CMD ["./subscriber","-is-adaptive=true", "-tunning=AMIGO","-execution-type=DynamicGoal", "-controller-type=SetpointWeighting", "-monitor-interval=5", "-prefetch-count=1",  "-max=1000", "-min=1", "-set-point=1000", "-direction=1",  "-kp=0.00101407", "-ki=0.00213378", "-kd=0.00012676", "-alfa=1","-beta=0.5"]
