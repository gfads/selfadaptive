rem Self generated file at 2023-08-03 18:15:31.7301033 -0300 -03 m=+0.044478001
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
set list=Dockerfile-BasicP-RootLocus Dockerfile-BasicP-Ziegler Dockerfile-BasicP-Cohen Dockerfile-BasicP-AMIGO Dockerfile-BasicPI-RootLocus Dockerfile-BasicPI-Ziegler Dockerfile-BasicPI-Cohen Dockerfile-BasicPI-AMIGO Dockerfile-BasicPID-RootLocus Dockerfile-BasicPID-Ziegler Dockerfile-BasicPID-Cohen Dockerfile-BasicPID-AMIGO Dockerfile-SmoothingDerivativePID-RootLocus Dockerfile-SmoothingDerivativePID-Ziegler Dockerfile-SmoothingDerivativePID-Cohen Dockerfile-SmoothingDerivativePID-AMIGO Dockerfile-IncrementalFormPID-RootLocus Dockerfile-IncrementalFormPID-Ziegler Dockerfile-IncrementalFormPID-Cohen Dockerfile-IncrementalFormPID-AMIGO Dockerfile-ErrorSquarePIDFull-RootLocus Dockerfile-ErrorSquarePIDFull-Ziegler Dockerfile-ErrorSquarePIDFull-Cohen Dockerfile-ErrorSquarePIDFull-AMIGO Dockerfile-ErrorSquarePIDProportional-RootLocus Dockerfile-ErrorSquarePIDProportional-Ziegler Dockerfile-ErrorSquarePIDProportional-Cohen Dockerfile-ErrorSquarePIDProportional-AMIGO Dockerfile-DeadZonePID-RootLocus Dockerfile-DeadZonePID-Ziegler Dockerfile-DeadZonePID-Cohen Dockerfile-DeadZonePID-AMIGO Dockerfile-GainScheduling-RootLocus Dockerfile-GainScheduling-Ziegler Dockerfile-GainScheduling-Cohen Dockerfile-GainScheduling-AMIGO Dockerfile-PIWithTwoDegreesOfFreedom-RootLocus Dockerfile-PIWithTwoDegreesOfFreedom-Ziegler Dockerfile-PIWithTwoDegreesOfFreedom-Cohen Dockerfile-PIWithTwoDegreesOfFreedom-AMIGO Dockerfile-WindUp-RootLocus Dockerfile-WindUp-Ziegler Dockerfile-WindUp-Cohen Dockerfile-WindUp-AMIGO Dockerfile-SetpointWeighting-RootLocus Dockerfile-SetpointWeighting-Ziegler Dockerfile-SetpointWeighting-Cohen Dockerfile-SetpointWeighting-AMIGO Dockerfile-AsTAR-RootLocus Dockerfile-AsTAR-Ziegler Dockerfile-AsTAR-Cohen Dockerfile-AsTAR-AMIGO Dockerfile-HPA-RootLocus Dockerfile-HPA-Ziegler Dockerfile-HPA-Cohen Dockerfile-HPA-AMIGO Dockerfile-OnOff-RootLocus Dockerfile-OnOff-Ziegler Dockerfile-OnOff-Cohen Dockerfile-OnOff-AMIGO Dockerfile-OnOffwithDeadZone-RootLocus Dockerfile-OnOffwithDeadZone-Ziegler Dockerfile-OnOffwithDeadZone-Cohen Dockerfile-OnOffwithDeadZone-AMIGO Dockerfile-OnOffwithHysteresis-RootLocus Dockerfile-OnOffwithHysteresis-Ziegler Dockerfile-OnOffwithHysteresis-Cohen Dockerfile-OnOffwithHysteresis-AMIGO 

echo ****** BEGIN OF EXPERIMENTS *******
for %x in (%!l(MISSING)ist%!)(MISSING) do (
echo %x 
   copy C:\Users\user\go\selfadaptive\temp\%x Dockerfile
   docker build --tag subscriber .
   docker run --rm --memory="1g" --cpus="1.0" -v C:\Users\user\go\selfadaptive\rabbitmq\data:/app/data subscriber
   del C:\Users\user\go\selfadaptive\temp\%x 
)
echo ****** END OF EXPERIMENTS *******
