#!/bin/bash

wget -O ui/electron.zip https://github.com/electron/electron/releases/download/v4.0.1/electron-v4.0.1-darwin-x64.zip
wget -O ui/astilectron.zip https://github.com/asticode/astilectron/archive/v0.34.0.zip
$GOPATH/bin/go-bindata ./ui/...
export GO111MODULE=on
export GOOS=darwin
export GOARCH=amd64
go build -mod=vendor
$GOPATH/bin/appify -name "Discord Bot GUI" -icon ./discord-512.png ./discord-bot-gui
zip -r discord-bot-gui_darwin.zip "Discord Bot GUI.app"