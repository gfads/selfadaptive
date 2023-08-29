@echo OFF
cls
set et=%1
set ct=%2
set t=%3
set f=%4

if "%~1"=="" goto :error1
if "%~2"=="" goto :error2
if "%~3"=="" goto :error3
if "%~4"=="" goto :error4

if %et% neq Experiment (
    if %et% neq ZieglerTraining (
        if %et% neq RootTraining (
              goto :error1
           )
    )
)
goto :ok

rem **** BEGIN OF OK ****
:ok
set GO111MODULE=on
set GOPATH=C:\Users\user\go;C:\Users\user\go\control\pkg\mod\github.com\streadway\amqp@v1.0.0;C:\Users\user\go\selfadaptive
set GOROOT=C:\Program Files\Go
set PATH=%PATH%;C:\Program Files\Docker\Docker

echo #### 1: Generate Dockerfiles/Batch file ####
cd C:\Users\user\go\selfadaptive\helper\gen
go run main.go -execution-type=%et% -controller-type=%ct% -tunning=%t% -output-file=%f%

cd C:\Users\user\go\selfadaptive

echo #### 2: Remove images ####
echo y | docker volume prune
echo y | docker image prune
echo y | docker container prune

echo #### 3: Execute Experiments ####
execute-all-experiments

echo #### 4: Generate Statistics ####
cd C:\Users\user\go\selfadaptive\helper\stats
rem main.exe

cd C:\Users\user\go\selfadaptive
goto :EOF
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