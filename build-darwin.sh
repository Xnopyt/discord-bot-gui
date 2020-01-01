#!/bin/zsh

$GOPATH/bin/go-bindata ./ui/...
export GO111MODULE=on
export GOOS=darwin
export GOARCH=amd64
go build -mod=vendor
$GOPATH/bin/appify -name "Discord Bot GUI" -icon ./discord-512.png ./discord-bot-gui
zip -r discord-bot-gui_darwin.zip "Discord Bot GUI.app"