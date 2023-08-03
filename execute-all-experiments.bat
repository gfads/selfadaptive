rem Self generated file at 2023-08-03 14:37:35.5351459 -0300 -03 m=+0.032399501
@echo off 
docker stop some-rabbit 
docker rm some-rabbit
docker run -d --memory="6g" --cpus="5.0" --name some-rabbit -p 5672:5672 rabbitmq
timeout /t 20
echo Removing previous containers
docker stop publisher
docker rm publisher
docker stop subscriber
docker rm subscriber
echo Removing previous volumes
echo y | docker volume prune
set list=Dockerfile-AsTAR-None Dockerfile-HPA-None Dockerfile-OnOff-None Dockerfile-OnOffwithDeadZone-None Dockerfile-OnOffwithHysteresis-None Dockerfile-BasicP-RootLocus Dockerfile-BasicP-Ziegler Dockerfile-BasicP-Cohen Dockerfile-BasicP-AMIGO Dockerfile-BasicPI-RootLocus Dockerfile-BasicPI-Ziegler Dockerfile-BasicPI-Cohen Dockerfile-BasicPI-AMIGO Dockerfile-BasicPID-RootLocus Dockerfile-BasicPID-Ziegler Dockerfile-BasicPID-Cohen Dockerfile-BasicPID-AMIGO Dockerfile-SmoothingDerivativePID-RootLocus Dockerfile-SmoothingDerivativePID-Ziegler Dockerfile-SmoothingDerivativePID-Cohen Dockerfile-SmoothingDerivativePID-AMIGO Dockerfile-IncrementalFormPID-RootLocus Dockerfile-IncrementalFormPID-Ziegler Dockerfile-IncrementalFormPID-Cohen Dockerfile-IncrementalFormPID-AMIGO Dockerfile-ErrorSquarePIDFull-RootLocus Dockerfile-ErrorSquarePIDFull-Ziegler Dockerfile-ErrorSquarePIDFull-Cohen Dockerfile-ErrorSquarePIDFull-AMIGO Dockerfile-ErrorSquarePIDProportional-RootLocus Dockerfile-ErrorSquarePIDProportional-Ziegler Dockerfile-ErrorSquarePIDProportional-Cohen Dockerfile-ErrorSquarePIDProportional-AMIGO Dockerfile-DeadZonePID-RootLocus Dockerfile-DeadZonePID-Ziegler Dockerfile-DeadZonePID-Cohen Dockerfile-DeadZonePID-AMIGO Dockerfile-GainScheduling-RootLocus Dockerfile-GainScheduling-Ziegler Dockerfile-GainScheduling-Cohen Dockerfile-GainScheduling-AMIGO Dockerfile-PIWithTwoDegreesOfFreedom-RootLocus Dockerfile-PIWithTwoDegreesOfFreedom-Ziegler Dockerfile-PIWithTwoDegreesOfFreedom-Cohen Dockerfile-PIWithTwoDegreesOfFreedom-AMIGO Dockerfile-WindUp-RootLocus Dockerfile-WindUp-Ziegler Dockerfile-WindUp-Cohen Dockerfile-WindUp-AMIGO Dockerfile-SetpointWeighting-RootLocus Dockerfile-SetpointWeighting-Ziegler Dockerfile-SetpointWeighting-Cohen Dockerfile-SetpointWeighting-AMIGO 

echo ****** BEGIN OF EXPERIMENTS *******
for %%x in (%list%) do (
echo %%x 
   copy C:\Users\user\go\selfadaptive\temp\%%x Dockerfile
   docker build --tag subscriber .
   docker run --rm --memory="1g" --cpus="1.0" -v C:\Users\user\go\selfadaptive\rabbitmq\data:/app/data subscriber
   del C:\Users\user\go\selfadaptive\temp\%%x 
)
echo ****** END OF EXPERIMENTS *******
 
