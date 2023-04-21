@echo off
rem docker stop some-rabbit
rem docker rm some-rabbit
rem docker run -d --memory="6g" --cpus="3.0" --name some-rabbit -p 5672:5672 rabbitmq
rem timeout /t 15

echo Create and Execute Publisher
copy Dockerfile-publisher Dockerfile
docker build --tag publisher .
START /B docker run publisher
rem START docker run publisher

echo Create and Execute Subscriber
copy Dockerfile-subscriber Dockerfile
docker build --tag subscriber .
START docker run subscriber

rem START docker run publisher