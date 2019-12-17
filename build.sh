#!/bin/bash

go-bindata ./ui/...
go build
./discord-bot-gui
