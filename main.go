package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

var wv *astilectron.Window

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go serveHTTP(ln)

	go evaulator()

	a, err := astilectron.New(astilectron.Options{AppName: "Discord Bot GUI"})
	if err != nil {
		astilog.Fatal(errors.Wrap(err, "main: creating astilectron failed"))
	}

	defer a.Close()

	if err = a.Start(); err != nil {
		astilog.Fatal(errors.Wrap(err, "main: starting astilectron failed"))
	}

	if wv, err = a.NewWindow("http://"+ln.Addr().String(), &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(720),
		Width:  astikit.IntPtr(1280),
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "main: new window failed"))
	}

	if err = wv.Create(); err != nil {
		astilog.Fatal(errors.Wrap(err, "main: creating window failed"))
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
			}
		}
		return nil
	})

	wv.OpenDevTools()

	a.Wait()
}
