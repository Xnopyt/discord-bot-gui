package main

import (
	"fmt"
	"html"
	"log"
	"net"
	"net/http"
	"runtime"
	"strings"

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
	route{"Login", "GET", "/", loginPage}, route{"Lbg", "GET", "/loginbg.jpg", lbg}, route{"DefAva", "GET", "/default.png", defaultavatar},
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
		if !strings.Contains(err.Error(), "use of closed network connection") {
			panic(err)
		}
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
		_, err := wv.Bind("bind", &mainBind{})
		if err != nil {
			log.Fatal(err)
		}
		wv.InjectCSS(string(MustAsset("ui/main.css")))
		if runtime.GOOS != "darwin" {
			wv.SetTitle("Discord Bot GUI - " + ses.State.User.String())
		}
		wv.Eval(string(MustAsset("ui/js/main.js")))
		wv.Eval(fmt.Sprintf(`
			document.getElementById("cname").innerHTML = %q;
			document.getElementById("cdiscriminator").innerHTML = '#%s';
			document.getElementById("cavatar").src = %q;
		`, html.EscapeString(ses.State.User.Username), ses.State.User.Discriminator, ses.State.User.AvatarURL("128")))
		loadServers()
		loadDMMembers()
	})
}

func ieUpdate() {
	browser.OpenURL("https://www.microsoft.com/en-us/download/internet-explorer.aspx")
}

func (m mainBind) Home() {
	currentServer = "HOME"
	currentChannel = ""
	wv.Dispatch(func() {
		wv.Eval(`loadhome()`)
		loadDMMembers()
	})
}
