#!/bin/bash

go-bindata ./ui/...
GO111MODULE=on go build -v -mod=vendor
./discord-bot-gui
