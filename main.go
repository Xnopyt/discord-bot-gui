package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/pkg/browser"
)

var wv *astilectron.Window

func main() {
	var l astikit.StdLogger
	electronzip := "ui/electron-" + runtime.GOOS + "_" + runtime.GOARCH + ".zip"
	electronProvisioner := astilectron.NewDisembedderProvisioner(Asset, "ui/astilectron.zip", electronzip, l)

	var ep string
	var err error
	if ep, err = os.Executable(); err != nil {
		log.Fatal(err)
	}
	vendorDir := filepath.Dir(ep)
	if vendorDir, err = filepath.Abs(vendorDir); err != nil {
		log.Fatal(err)
	}
	if v := os.Getenv("APPDATA"); len(v) > 0 {
		vendorDir = filepath.Join(v, "Discord Bot GUI")
		return
	}
	vendorDir = filepath.Join(vendorDir, "vendor")
	if _, err := os.Stat(vendorDir); os.IsNotExist(err) {
		os.Mkdir(vendorDir, os.ModePerm)
	}
	ico := filepath.Join(vendorDir, "dbg.ico")
	icns := filepath.Join(vendorDir, "dbg.icns")
	var icnsOut *os.File
	var icoOut *os.File
	if icnsOut, err = os.Create(icns); err != nil {
		log.Fatal(err)
	}
	if icoOut, err = os.Create(ico); err != nil {
		log.Fatal(err)
	}
	_, err = icnsOut.Write(MustAsset("ui/assets/discord-512.icns"))
	if err != nil {
		icnsOut.Close()
		log.Fatal(err)
	}
	icnsOut.Close()
	_, err = icoOut.Write(MustAsset("ui/assets/discord-512.ico"))
	if err != nil {
		icoOut.Close()
		log.Fatal(err)
	}
	icoOut.Close()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go serveHTTP(ln)

	go evaulator()

	a, err := astilectron.New(l, astilectron.Options{
		AppName:            "Discord Bot GUI",
		AppIconDarwinPath:  icns,
		AppIconDefaultPath: ico,
	})
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

			case "openURL":
				if msg.Content != "" {
					browser.OpenURL(msg.Content)
				}

			case "sendFile":
				sendFile(msg.Content)
			}
		}
		return nil
	})

	a.Wait()
}
