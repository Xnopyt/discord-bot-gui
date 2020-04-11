package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/pkg/browser"
	"github.com/zserge/webview"
)

var wv webview.WebView

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go serveHTTP(ln)
	wv = webview.New(true)
	defer wv.Destroy()
	wv.SetTitle("Discord Bot GUI - Login")
	wv.SetSize(1280, 720, webview.HintNone)
	wv.Bind("wv", webviewCallback)
	wv.Navigate("http://" + ln.Addr().String())
	wv.Run()
}

func webviewCallback(s string) {
	callback, ok := wvCallbacks[s]
	if ok {
		callback()
	} else {
		var msg uiMsg
		err := json.Unmarshal([]byte(s), &msg)
		if err != nil {
			return
		}
		switch msg.Type {
		case "connect":
			connect(msg.Content)

		case "selectTargetServer":
			selectTargetServer(msg.Content)

		case "setActiveChannel":
			setActiveChannel(msg.Content)

		case "sendMessage":
			sendMessage(msg.Content)

		case "loadDMChannel":
			loadDMChannel(msg.Content)

		case "openURL":
			if msg.Content != "" {
				browser.OpenURL(msg.Content)
			}

		case "sendFile":
			sendFile(msg.Content)
		}
	}
}
