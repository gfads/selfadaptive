@echo off
set GO111MODULE=on
set GOPATH=C:\Users\user\go;C:\Users\user\go\control\pkg\mod\github.com\streadway\amqp@v1.0.0;C:\Users\user\go\selfadaptive
set GOROOT=C:\Program Files\Go

rem start rabbitmq
docker stop some-rabbit
docker rm some-rabbit
docker run -d --memory="6g" --cpus="5.0" --name some-rabbit -p 5672:5672 rabbitmq

timeout /t 20

rem subscriber
docker stop some-subscriber
docker rm some-subscriber
copy C:\Users\user\go\selfadaptive\Dockerfile-subscriber Dockerfile
docker build --tag subscriber .
docker run -d --rm --name some-subscriber --memory="1g" --cpus="1.0" -v C:\Users\user\go\selfadaptive\rabbitmq\data:/app/data subscriber

rem publisher
rem docker stop some-publisher
rem docker rm some-publisher
rem copy C:\Users\user\go\selfadaptive\Dockerfile-publisher Dockerfile
rem docker build --tag publisher .
rem docker run --rm --name some-publisher --memory="1g" --cpus="1.0" publisher
