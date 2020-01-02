#!/bin/bash

go-bindata ./ui/...
xgo -targets darwin-10.11/amd64 -branch webview github.com/Xnopyt/discord-bot-gui
appify -name "Discord Bot GUI" -icon ./discord-512.png ./discord-bot-gui-darwin-10.11-amd64
zip -r discord-bot-gui_darwin.zip "Discord Bot GUI.app"
rm -rf "Discord Bot GUI.app" discord-bot-gui-darwin-10.11-amd64