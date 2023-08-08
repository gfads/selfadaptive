rem This file has been generated automatically at 2023-08-08 08:19:48.3477579 -0300 -03 m=+0.068799801
@echo off 
docker stop some-rabbit 
docker rm some-rabbit
docker run -d --memory="6g" --cpus="5.0" --name some-rabbit -p 5672:5672 rabbitmq
timeout /t 10
docker stop publisher
docker rm publisher
docker stop subscriber
docker rm subscriber
set list=Dockerfile-BasicPI-RootLocus Dockerfile-BasicPI-AMIGO Dockerfile-BasicPID-RootLocus Dockerfile-BasicPID-AMIGO Dockerfile-SmoothingDerivativePID-RootLocus Dockerfile-SmoothingDerivativePID-AMIGO Dockerfile-IncrementalFormPID-RootLocus Dockerfile-IncrementalFormPID-AMIGO Dockerfile-ErrorSquarePIDFull-RootLocus Dockerfile-ErrorSquarePIDFull-AMIGO Dockerfile-ErrorSquarePIDProportional-RootLocus Dockerfile-ErrorSquarePIDProportional-AMIGO Dockerfile-DeadZonePID-RootLocus Dockerfile-DeadZonePID-AMIGO Dockerfile-GainScheduling-RootLocus Dockerfile-GainScheduling-AMIGO Dockerfile-PIWithTwoDegreesOfFreedom-RootLocus Dockerfile-PIWithTwoDegreesOfFreedom-AMIGO Dockerfile-WindUp-RootLocus Dockerfile-WindUp-AMIGO Dockerfile-SetpointWeighting-RootLocus Dockerfile-SetpointWeighting-AMIGO Dockerfile-AsTAR-None Dockerfile-AsTAR-None Dockerfile-HPA-None Dockerfile-HPA-None 

for %%x in (%list%) do (
echo %%x 
   copy C:\Users\user\go\selfadaptive\temp\%%x Dockerfile
   docker build --tag subscriber .
   docker run --rm --name some-subscriber --memory="1g" --cpus="1.0" -v C:\Users\user\go\selfadaptive\rabbitmq\data:/app/data subscriber
   del C:\Users\user\go\selfadaptive\temp\%%x 
   echo y | docker volume prune 
   echo y | docker image prune 
)
