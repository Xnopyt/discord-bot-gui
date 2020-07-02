package main

import (
	"encoding/json"

	"github.com/pkg/browser"
	"github.com/zserge/webview"
)

var wv webview.WebView

func main() {
	wv = webview.New(true)
	defer wv.Destroy()
	wv.SetTitle("Discord Bot GUI - Login")
	wv.SetSize(1280, 720, webview.HintNone)
	wv.Bind("wv", webviewCallback)
	wv.Bind("readClipboard", readClipboard)
	wv.Bind("writeClipboard", writeClipboard)
	wv.Bind("deleteMessage", deleteMessage)
	wv.Navigate("data:text/html," + string(MustAsset("ui/login.html")))
	wv.Run()
	if ses != nil {
		ses.Close()
	}
}

func webviewCallback(s string) {
	callback, ok := wvCallbacks[s]
	if ok {
		go callback()
	} else {
		var msg uiMsg
		err := json.Unmarshal([]byte(s), &msg)
		if err != nil {
			return
		}
		switch msg.Type {
		case "connect":
			go connect(msg.Content)

		case "selectTargetServer":
			go selectTargetServer(msg.Content)

		case "setActiveChannel":
			go setActiveChannel(msg.Content)

		case "sendMessage":
			go sendMessage(msg.Content)

		case "loadDMChannel":
			go loadDMChannel(msg.Content)

		case "openURL":
			if msg.Content != "" {
				go browser.OpenURL(msg.Content)
			}

		case "sendFile":
			go sendFile(msg.Content)
		}
	}
}
