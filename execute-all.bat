rem Self generated file at 2023-08-01 11:31:35.562085 -0300 -03 m=+0.014342801
@echo off 
docker stop some-rabbit 
docker rm some-rabbit
docker run -d --memory="6g" --cpus="5.0" --name some-rabbit -p 5672:5672 rabbitmq
timeout /t 10
@echo Removing previous containers
docker stop publisher
docker rm publisher
docker stop subscriber
docker rm subscriber
set list=Dockerfile-PIWithTwoDegreesOfFreedom-RootLocus Dockerfile-PIWithTwoDegreesOfFreedom-Ziegler Dockerfile-PIWithTwoDegreesOfFreedom-Cohen Dockerfile-PIWithTwoDegreesOfFreedom-AMIGO 
echo ****** BEGIN OF EXPERIMENTS *******
for %%x in (%list%) do (
   copy %%x Dockerfile
   docker build --tag subscriber .
	docker run --rm --memory="1g" --cpus="1.0" -v C:\Users\user\go\selfadaptive\rabbitmq\data:/app/data subscriber
   del %%x 
)
echo ****** END OF EXPERIMENTS *******
 
