#! /bin/sh

export GOPROXY=https://goproxy.cn
go mod tidy
go run main.go
