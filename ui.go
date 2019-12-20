package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/browser"
	"github.com/zserge/webview"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []route{
	route{"Login", "GET", "/", loginPage}, route{"Lbg", "GET", "/loginbg.jpg", lbg}, route{"Main", "GET", "/main", mainPage}, route{"DefAva", "GET", "/default.png", defaultavatar},
}

var wvCallbacks map[string]func()

func init() {
	wvCallbacks = make(map[string]func())

	wvCallbacks["loginSetup"] = loginSetup
	wvCallbacks["mainSetup"] = mainSetup
	wvCallbacks["ieUpdate"] = ieUpdate
}

func newRouter() *mux.Router {
	router := mux.NewRouter()
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

func loginPage(rw http.ResponseWriter, r *http.Request) {
	_, err := rw.Write(MustAsset("ui/login.html"))
	if err != nil {
		log.Fatal(err)
	}
}

func mainPage(rw http.ResponseWriter, r *http.Request) {
	_, err := rw.Write(MustAsset("ui/main.html"))
	if err != nil {
		log.Fatal(err)
	}
}

func lbg(rw http.ResponseWriter, r *http.Request) {
	_, err := rw.Write(MustAsset("ui/assets/loginbg.jpg"))
	if err != nil {
		log.Fatal(err)
	}
}

func defaultavatar(rw http.ResponseWriter, r *http.Request) {
	_, err := rw.Write(MustAsset("ui/assets/default-avatar.png"))
	if err != nil {
		log.Fatal(err)
	}
}

func serveHTTP(ln net.Listener) {
	router := newRouter()
	if err := http.Serve(ln, router); err != nil {
		fmt.Println(err)
	}
}

func webviewCallback(w webview.WebView, s string) {
	callback, ok := wvCallbacks[s]
	if ok {
		callback()
	} else {
		fmt.Println("Attempted to call unknown function " + s)
	}
}

func loginSetup() {
	wv.Dispatch(func() {
		_, err := wv.Bind("binder", &binder{})
		if err != nil {
			log.Fatal(err)
		}
		err = wv.Eval(string(MustAsset("ui/js/login.js")))
		if err != nil {
			log.Fatal(err)
		}
		wv.InjectCSS(string(MustAsset("ui/login.css")))
	})
}

func mainSetup() {
	wv.Dispatch(func() {
		wv.InjectCSS(string(MustAsset("ui/main.css")))
		wv.SetTitle("Discord Bot GUI")
	})
}

func ieUpdate() {
	browser.OpenURL("https://www.microsoft.com/en-us/download/internet-explorer.aspx")
}
