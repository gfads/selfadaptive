@echo off
docker stop some-rabbit
docker rm some-rabbit
docker run -d --memory="6g" --cpus="3.0" --name some-rabbit -p 5672:5672 rabbitmq
timeout /t 30

docker stop publisher
docker rm publisher

docker stop subscriber
docker rm subscriber

echo Create and Execute Publisher
copy Dockerfile-publisher Dockerfile
docker build --tag publisher .

set clients=1
:loop
  START /B docker run publisher
  set /a clients=clients-1
  if %clients%==0 goto exitloop
  goto loop
:exitloop

echo Create and Execute Subscriber
copy Dockerfile-subscriber Dockerfile
docker build --tag subscriber .
START docker run subscriber

rem START docker run publisher