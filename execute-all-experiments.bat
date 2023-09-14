rem This file has been generated automatically at 2023-09-14 13:39:42.1952815 -0300 -03 m=+0.014782901
@echo off 
timeout /t 10
docker stop publisher
docker rm publisher
docker stop subscriber
docker rm subscriber
set list=Dockerfile-Experiment-IncrementalFormPID-AMIGO 

for %%x in (%list%) do (
echo %%x 
   copy C:\Users\user\go\selfadaptive\temp\%%x Dockerfile
   docker build --tag subscriber .
   docker run --rm --name some-subscriber --memory="1g" --cpus="1.0" -v C:\Users\user\go\selfadaptive\rabbitmq\data:/app/data subscriber
   del C:\Users\user\go\selfadaptive\temp\%%x 
   echo y| docker volume prune 
   echo y| docker image prune 
)
