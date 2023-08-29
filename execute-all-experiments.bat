rem This file has been generated automatically at 2023-08-29 14:01:28.3476293 -0300 -03 m=+0.013409001
@echo off 
docker stop some-rabbit 
docker rm some-rabbit
docker run -d --memory="6g" --cpus="5.0" --name some-rabbit -p 5672:5672 rabbitmq
timeout /t 10
docker stop publisher
docker rm publisher
docker stop subscriber
docker rm subscriber
set list=Dockerfile-Experiment-BasicPID-Ziegler 

for %%x in (%list%) do (
echo %%x 
   copy C:\Users\user\go\selfadaptive\temp\%%x Dockerfile
   docker build --tag subscriber .
   docker run --rm --name some-subscriber --memory="1g" --cpus="1.0" -v C:\Users\user\go\selfadaptive\rabbitmq\data:/app/data subscriber
   del C:\Users\user\go\selfadaptive\temp\%%x 
   echo y | docker volume prune 
   echo y | docker image prune 
)
docker stop some-rabbit 
docker rm some-rabbit
