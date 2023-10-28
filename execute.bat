@echo OFF
cls

rem remove previous rabbitmq
docker stop some-rabbit
docker rm some-rabbit

rem execute new instance of rabbitmq
docker run -d --memory=6g --cpus=5.0 --name some-rabbit -p 5672:5672 rabbitmq

rem configure variables
rem set et=Experiment
set et=ExperimentalDesign
set ct=BasicPID
set t=RootLocus
set b=%et%-%ct%-%t%
rem set f=1 2 3 4 5 6 7 8 9 10
set f=1

for %%x in (%f%) do (
    set GO111MODULE=on
    set GOPATH=C:\Users\user\go;C:\Users\user\go\control\pkg\mod\github.com\streadway\amqp@v1.0.0;C:\Users\user\go\selfadaptive
    set GOROOT=C:\Program Files\Go
    set PATH=%PATH%;C:\Program Files\Docker\Docker

    echo #### 1: Generate Dockerfiles/Batch file ####
    cd C:\Users\user\go\selfadaptive\helper\gen
    go run main.go -execution-type=%et% -controller-type=%ct% -tunning=%t% -output-file=%b%-%%x

    cd C:\Users\user\go\selfadaptive

    echo #### 2: Remove images ####
    echo y| docker volume prune
    echo y| docker image prune
    echo y| docker container prune

    echo #### 3: Execute Experiments ####
    execute-all-experiments
)

echo #### 4: Clear data & Stop RabbitMQ ####
echo y| docker volume prune
echo y| docker image prune
echo y| docker container prune

docker stop some-rabbit
docker rm some-rabbit

echo #### 4: Generate Statistics ####
cd C:\Users\user\go\selfadaptive\helper\stats
rem main.exe

cd C:\Users\user\go\selfadaptive

rem goto :EOF
rem **** END OF OK ****

:error1
echo ERROR:: Execution type is invalid. Use any of the following options: Experiment, ZiegletTraining, RootTraining
goto :EOF
rem **** END OF ERROR ****

:error2
echo ERROR:: Controller type is invalid. Use any of the following options: xxx, yyy, zzz
goto :EOF
rem **** END OF ERROR ****

:error3
echo ERROR:: Tuning type is invalid. Use any of the following options: xxx, yyy, zzz
goto :EOF
rem **** END OF ERROR ****

:error4
echo ERROR:: File name is invalid.
goto :EOF
rem **** END OF ERROR ****