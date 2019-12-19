#!/bin/bash

go-bindata ./ui/...
GO111MODULE=on go build -mod=vendor
./discord-bot-gui
