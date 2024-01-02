#! /bin/sh

export GOPROXY=https://goproxy.cn
go mod tidy
go install github.com/go-delve/delve/cmd/dlv@latest
/go/bin/dlv --headless --log --listen :2345 --api-version 2 --accept-multiclient debug main.go
