@echo OFF
cls

rem remove previous rabbitmq
docker stop myrabbit
docker rm myrabbit

rem execute new instance of rabbitmq
rem docker run -d --memory=6g --cpus=5.0 --name myrabbit -p 5672:5672 rabbitmq
docker run -d --memory=6g --cpus=5.0 --name myrabbit -p 5672:5672 rabbitmq

rem Placed here to fix the error "\Common inesperado" TODO
rem PATH=C:\Program Files (x86)\Common Files\Oracle\Java\javapath;C:\WINDOWS\system32;C:\WINDOWS;C:\WINDOWS\System32\Wbem;C:\WINDOWS\System32\WindowsPowerShell\v1.0\;C:\WINDOWS\System32\OpenSSH\;C:\Program Files\Go\bin;C:\Program Files\Git\cmd;C:\Program Files\Docker\Docker\resources\bin;C:\Program Files\MATLAB\R2023b\bin;C:\Users\user\AppData\Local\Microsoft\WindowsApps;C:\Program Files\JetBrains\GoLand 2022.1.2\bin;;C:\Users\user\go\bin;C:\MinGW\bin;C:\Users\user\AppData\Local\Google\Cloud SDK\google-cloud-sdk\bin;C:\Program Files\Docker\Docker
set PATH=C:\WINDOWS\system32;C:\WINDOWS;C:\WINDOWS\System32\Wbem;C:\WINDOWS\System32\WindowsPowerShell\v1.0\;C:\WINDOWS\System32\OpenSSH\;C:\Program Files\Go\bin;C:\Program Files\Git\cmd;C:\Program Files\Docker\Docker\resources\bin;C:\Program Files\MATLAB\R2023b\bin;C:\Users\user\AppData\Local\Microsoft\WindowsApps;C:\Program Files\JetBrains\GoLand 2022.1.2\bin;;C:\Users\user\go\bin;C:\MinGW\bin;C:\Users\user\AppData\Local\Google\Cloud SDK\google-cloud-sdk\bin;C:\Program Files\Docker\Docker
set PATH=%PATH%;C:\Program Files\Docker\Docker

rem configure variables
rem For ExperimentalDesign (ramp/sine training), use ct=BasicPID and t=None
rem set et=ExperimentalDesign
set et=Experiment
rem set et=RootTraining
rem set et=ZieglerTraining
rem set ct=ErrorSquarePIDFull
rem set ct=ErrorSquarePIDProportional
rem set ct=BasicPI
rem set ct=GainScheduling
rem set ct=HPA
rem set ct=Fuzzy
rem set ct=IncrementalFormPID
set ct=AsTAR
rem set t=RootLocus
set t=None
rem set t=Ziegler
rem set t=AMIGO
set b=%et%-%ct%-%t
rem set f=1 2 3 4 5 6 7 8 9 10
set f=2

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

docker stop myrabbit
docker rm myrabbit

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