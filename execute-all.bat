rem Self generated file at 2023-08-02 05:59:31.7814977 -0300 -03 m=+0.057063401
@echo off 
docker stop some-rabbit 
docker rm some-rabbit
docker run -d --memory="6g" --cpus="5.0" --name some-rabbit -p 5672:5672 rabbitmq
timeout /t 20
@echo Removing previous containers
docker stop publisher
docker rm publisher
docker stop subscriber
docker rm subscriber
@echo Removing previous volumes
echo y | docker volume prune
set list=Dockerfile-AsTAR-None Dockerfile-HPA-None Dockerfile-BasicP-RootLocus Dockerfile-BasicP-Ziegler Dockerfile-BasicP-Cohen Dockerfile-BasicPI-RootLocus Dockerfile-BasicPI-Ziegler Dockerfile-BasicPI-Cohen Dockerfile-BasicPI-AMIGO Dockerfile-BasicPID-RootLocus Dockerfile-BasicPID-Ziegler Dockerfile-BasicPID-Cohen Dockerfile-BasicPID-AMIGO Dockerfile-DeadZonePID-RootLocus Dockerfile-DeadZonePID-Ziegler Dockerfile-DeadZonePID-Cohen Dockerfile-DeadZonePID-AMIGO Dockerfile-IncrementalFormPID-RootLocus Dockerfile-IncrementalFormPID-Ziegler Dockerfile-IncrementalFormPID-Cohen Dockerfile-IncrementalFormPID-AMIGO Dockerfile-SmoothingDerivativePID-RootLocus Dockerfile-SmoothingDerivativePID-Ziegler Dockerfile-SmoothingDerivativePID-Cohen Dockerfile-SmoothingDerivativePID-AMIGO Dockerfile-GainScheduling-RootLocus Dockerfile-GainScheduling-Ziegler Dockerfile-GainScheduling-Cohen Dockerfile-GainScheduling-AMIGO Dockerfile-SetpointWeighting-RootLocus Dockerfile-SetpointWeighting-Ziegler Dockerfile-SetpointWeighting-Cohen Dockerfile-SetpointWeighting-AMIGO Dockerfile-SetpointWeighting-RootLocus Dockerfile-SetpointWeighting-Ziegler Dockerfile-SetpointWeighting-Cohen Dockerfile-SetpointWeighting-AMIGO Dockerfile-ErrorSquarePIDProportional-RootLocus Dockerfile-ErrorSquarePIDProportional-Ziegler Dockerfile-ErrorSquarePIDProportional-Cohen Dockerfile-ErrorSquarePIDProportional-AMIGO Dockerfile-ErrorSquarePIDFull-RootLocus Dockerfile-ErrorSquarePIDFull-Ziegler Dockerfile-ErrorSquarePIDFull-Cohen Dockerfile-ErrorSquarePIDFull-AMIGO 
echo ****** BEGIN OF EXPERIMENTS *******
for %%x in (%list%) do (
   copy %%x Dockerfile
   docker build --tag subscriber .
	docker run --rm --memory="1g" --cpus="1.0" -v C:\Users\user\go\selfadaptive\rabbitmq\data:/app/data subscriber
   del %%x 
)
echo ****** END OF EXPERIMENTS *******
 
