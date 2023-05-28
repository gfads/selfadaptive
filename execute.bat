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

echo ****** Create and Execute Publisher (Local publisher) ******
copy Dockerfile-publisher Dockerfile
docker build --tag publisher .
START /B docker run publisher *** unrem if publisher executes in this host

rem set clients=1
rem :loop
rem  START /B docker run publisher
rem  set /a clients=clients-1
rem  if %clients%==0 goto exitloop
rem  goto loop
rem :exitloop

timeout /t 10

echo ****** Create and Execute Subscriber ******
rem copy Dockerfile-subscriber Dockerfile
rem docker build --tag subscriber .
rem docker run subscriber