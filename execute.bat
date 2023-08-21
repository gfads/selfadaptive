@echo off
set GO111MODULE=on
set GOPATH=C:\Users\user\go;C:\Users\user\go\control\pkg\mod\github.com\streadway\amqp@v1.0.0;C:\Users\user\go\selfadaptive
set GOROOT=C:\Program Files\Go

echo #### Generate Dockerfiles/Batch file ####
cd C:\Users\user\go\selfadaptive\helper\gen
rem go run main.go -execution-type=Experiment
rem go run main.go -execution-type=Static
rem go run main.go -execution-type=RootLocusTraining
go run main.go -execution-type=ZieglerTraining

cd C:\Users\user\go\selfadaptive

echo ##### Remove images ####
echo y | docker volume prune
echo y | docker image prune
echo y | docker container prune

echo #### Execute Experiments ####
execute-all-experiments

echo #### Generate Statistics ####
cd C:\Users\user\go\selfadaptive\helper\stats
main.exe

cd C:\Users\user\go\selfadaptive
