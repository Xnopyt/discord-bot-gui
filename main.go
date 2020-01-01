package main

import (
	"log"
	"net"
	"runtime"

	"github.com/zserge/webview"
)

var wv webview.WebView

func main() {
	if runtime.GOOS == "windows" {
		log.Fatal("Discord Bot GUI only currently supports linux and macOS")
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go serveHTTP(ln)
	var title string
	if runtime.GOOS == "darwin" {
		title = "Discord Bot GUI"
	} else {
		title = "Discord Bot GUI - Login"
	}
	wv = webview.New(webview.Settings{
		Title:                  title,
		URL:                    "http://" + ln.Addr().String(),
		Width:                  1280,
		Height:                 720,
		Resizable:              true,
		Debug:                  true,
		ExternalInvokeCallback: webviewCallback,
	})
	wv.Run()
}
