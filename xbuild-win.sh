#!/bin/bash

go-bindata ./ui/...
export GO111MODULE=on
export GOOS=windows
export GOARCH=amd64
rsrc -ico=discord-512.ico -arch="$GOARCH" -o=discord-bot-gui.syso
go build -v -mod=vendor
rm discord-bot-gui.syso