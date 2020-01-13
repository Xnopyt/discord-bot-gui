#!/bin/bash

wget -O ui/electron.zip https://github.com/electron/electron/releases/download/v4.0.1/electron-v4.0.1-linux-x64.zip
wget -O ui/astilectron.zip https://github.com/asticode/astilectron/archive/v0.34.0.zip
go-bindata ./ui/...
export GOARCH=amd64
export GOOS=linux
GO111MODULE=on go build -v -mod=vendor
./discord-bot-gui
