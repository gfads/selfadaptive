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

echo ****** Create and Execute Subscriber 1 ******
copy Dockerfile-sub-setpoint-cohen Dockerfile
docker build --tag subscriber .
docker run --memory="1g" --cpus="1.0" -v C:\Users\user\go\selfadaptive\rabbitmq\data:/app/data subscriber

echo ****** Create and Execute Subscriber 2 ******
copy Dockerfile-sub-setpoint-amigo Dockerfile
docker build --tag subscriber .
docker run --memory="1g" --cpus="1.0" -v C:\Users\user\go\selfadaptive\rabbitmq\data:/app/data subscriber
