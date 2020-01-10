package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
)

var wv *astilectron.Window

func main() {
	electronProvisioner := astilectron.NewDisembedderProvisioner(Asset, "ui/astilectron.zip", "ui/electron.zip")
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go serveHTTP(ln)

	go evaulator()

	a, err := astilectron.New(astilectron.Options{AppName: "Discord Bot GUI"})
	if err != nil {
		log.Fatal(err)
	}

	defer a.Close()

	a.SetProvisioner(electronProvisioner)

	if err = a.Start(); err != nil {
		log.Fatal(err)
	}

	if wv, err = a.NewWindow("http://"+ln.Addr().String(), &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(720),
		Width:  astikit.IntPtr(1280),
	}); err != nil {
		log.Fatal(err)
	}

	if err = wv.Create(); err != nil {
		log.Fatal(err)
	}

	wv.OnMessage(func(m *astilectron.EventMessage) interface{} {
		var s string
		m.Unmarshal(&s)

		callback, ok := wvCallbacks[s]
		if ok {
			callback()
		} else {
			var msg uiMsg
			err := json.Unmarshal([]byte(s), &msg)
			if err != nil {
				return nil
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
			}
		}
		return nil
	})

	//wv.OpenDevTools()

	a.Wait()
}
