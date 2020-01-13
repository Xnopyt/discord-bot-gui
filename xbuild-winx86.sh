#!/bin/bash

wget -O ui/electron.zip https://github.com/electron/electron/releases/download/v4.0.1/electron-v4.0.1-win32-ia32.zip
wget -O ui/astilectron.zip https://github.com/asticode/astilectron/archive/v0.34.0.zip
go-bindata ./ui/...
export GO111MODULE=on
export GOOS=windows
export GOARCH=386
rsrc -ico=discord-512.ico -arch="$GOARCH" -o=discord-bot-gui.syso
go build -v -mod=vendor
rm discord-bot-gui.syso