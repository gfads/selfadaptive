rem This file has been generated automatically at 2023-08-04 17:20:38.0521398 -0300 -03 m=+0.027032601
@echo off 
docker stop some-rabbit 
docker rm some-rabbit
docker run -d --memory="6g" --cpus="5.0" --name some-rabbit -p 5672:5672 rabbitmq
timeout /t 10
docker stop publisher
docker rm publisher
docker stop subscriber
docker rm subscriber
set list=Dockerfile-BasicPID-RootLocus Dockerfile-SmoothingDerivativePID-RootLocus Dockerfile-IncrementalFormPID-RootLocus Dockerfile-ErrorSquarePIDFull-RootLocus Dockerfile-ErrorSquarePIDProportional-RootLocus Dockerfile-DeadZonePID-RootLocus Dockerfile-GainScheduling-RootLocus Dockerfile-PIWithTwoDegreesOfFreedom-RootLocus Dockerfile-WindUp-RootLocus Dockerfile-SetpointWeighting-RootLocus Dockerfile-AsTAR-None Dockerfile-HPA-None 

for %%x in (%list%) do (
echo %%x 
   copy C:\Users\user\go\selfadaptive\temp\%%x Dockerfile
   docker build --tag subscriber .
   docker run --rm --name some-subscriber --memory="1g" --cpus="1.0" -v C:\Users\user\go\selfadaptive\rabbitmq\data:/app/data subscriber
   del C:\Users\user\go\selfadaptive\temp\%%x 
)
