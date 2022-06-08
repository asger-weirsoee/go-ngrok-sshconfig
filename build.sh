#!/bin/sh

GOOS=linux GOARCH=amd64 go build -o bin/gogrok-amd64-linux main.go
GOOS=linux GOARCH=386 go build -o bin/gogrok-386-linux main.go

GOOS=windows GOARCH=386 go build -o bin/gogrok-386.exe main.go
GOOS=windows GOARCH=amd64 go build -o bin/gogrok-amd64.exe main.go