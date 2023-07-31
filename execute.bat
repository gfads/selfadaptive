@echo off
docker stop some-rabbit
docker rm some-rabbit
rem docker run -d --memory="6g" --cpus="3.0" --name some-rabbit -p 5672:5672 rabbitmq
docker run -d --memory="6g" --cpus="5.0" --name some-rabbit -p 5672:5672 rabbitmq

timeout /t 10

@echo Removing previous containers
docker stop publisher
docker rm publisher

docker stop subscriber
docker rm subscriber

rem echo Execute Dockerfile initialiser
rem TODO
rem copy Docker-initialiser Dockerfile

rem set list=Dockerfile-sub-astar Dockerfile-sub-hpa Dockerfile-sub-setpoint-ziegler Dockerfile-sub-setpoint-cohen Dockerfile-sub-setpoint-root
rem set list=Dockerfile-sub-gain-error
rem set list=Dockerfile-sub-p-root Dockerfile-sub-p-ziegler Dockerfile-sub-p-cohen
rem set list=Dockerfile-sub-onoff-basic Dockerfile-sub-onoff-deadzone Dockerfile-sub-onoff-hysteresis
rem set list=Dockerfile-sub-two-root Dockerfile-sub-two-ziegler Dockerfile-sub-two-cohen Dockerfile-sub-two-amigo
rem set list=Dockerfile-sub-two-cohen Dockerfile-sub-two-amigo
rem set list=Dockerfile-subscriber
rem set list=Dockerfile-sub-astar
rem set list=Dockerfile-sub-astar Dockerfile-sub-hpa Dockerfile-sub-onoff-basic Dockerfile-sub-onoff-deadzone Dockerfile-sub-hysteresis Dockerfile-sub-p-root Dockerfile-sub-p-zigler Dockerfile-sub-p-cohen Dockerfile-sub-pi-root Dockerfile-sub-pi-zigler Dockerfile-sub-pi-cohen Dockerfile-sub-pi-amigo Dockerfile-sub-pid-root Dockerfile-sub-pid-zigler Dockerfile-sub-pid-cohen Dockerfile-sub-pid-amigo
rem set list=Dockerfile-sub-astar Dockerfile-sub-hpa Dockerfile-sub-p-root Dockerfile-sub-p-ziegler Dockerfile-sub-p-cohen Dockerfile-sub-pi-root Dockerfile-sub-pi-zigler Dockerfile-sub-pi-cohen Dockerfile-sub-pi-amigo Dockerfile-sub-pid-root Dockerfile-sub-pid-zigler Dockerfile-sub-pid-cohen Dockerfile-sub-pid-amigo
set list=Apague

echo ****** BEGIN OF EXPERIMENTS *******
for %%x in (%list%) do (
        echo %%x
        copy %%x Dockerfile
        docker build --tag subscriber .
        docker run --memory="1g" --cpus="1.0" -v C:\Users\user\go\selfadaptive\rabbitmq\data:/app/data subscriber
       )
echo ****** END OF EXPERIMENTS *******