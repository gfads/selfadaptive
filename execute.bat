@echo off
set GO111MODULE=on
set GOPATH=C:\Users\user\go;C:\Users\user\go\control\pkg\mod\github.com\streadway\amqp@v1.0.0;C:\Users\user\go\selfadaptive
set GOROOT=C:\Program Files\Go

c:
cd C:\Users\user\go\selfadaptive\gen
go build -o main.exe main.go
main.exe
cd C:\Users\user\go\selfadaptive
execute-all
