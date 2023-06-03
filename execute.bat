@echo off
docker stop some-rabbit
docker rm some-rabbit
docker run -d --memory="6g" --cpus="3.0" --name some-rabbit -p 5672:5672 rabbitmq
timeout /t 10

@echo Removing previous containers
docker stop publisher
docker rm publisher

docker stop subscriber
docker rm subscriber

rem set clients=1
rem :loop
rem  START /B docker run publisher
rem  set /a clients=clients-1
rem  if %clients%==0 goto exitloop
rem  goto loop
rem :exitloop

timeout /t 10

echo ****** Create and Execute Subscriber ******
copy Dockerfile-subscriber Dockerfile
docker build --tag subscriber .
rem START /B docker run --memory="2g" --cpus="2.0" subscriber
rem START docker run --memory="6g" --cpus="3.0" subscriber
docker run --memory="6g" --cpus="3.0" subscriber

echo ****** Create and Execute Publisher (Local publisher) ******
rem copy Dockerfile-publisher Dockerfile
rem docker build --tag publisher .
rem START /B docker run publisher
rem docker run publisher
