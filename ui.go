package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
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
		wv.InjectCSS(string(MustAsset("ui/main.css")))
		wv.SetTitle("Discord Bot GUI - " + ses.State.User.String())
		wv.Eval(`
			document.getElementById("cname").innerHTML = '` + ses.State.User.Username + `';
			document.getElementById("cdiscriminator").innerHTML = '#` + ses.State.User.Discriminator + `';
			document.getElementById("cavatar").src = '` + ses.State.User.AvatarURL("128") + `';
		`)
		guilds, err := ses.UserGuilds(100, "", "")
		if err != nil {
			panic(err)
		}
		for _, v := range guilds {
			guild, _ := ses.Guild(v.ID)
			if guild.IconURL() == "" {
				var shortname string
				words := strings.Split(guild.Name, " ")
				for _, word := range words {
					if len(shortname) > 4 {
						break
					}
					shortname += string(word[0])
				}
				wv.Eval(`
					var newserver = document.createElement("div");
					newserver.className = "server";
					newserver.id = "` + guild.ID + `";
					var newsel = document.createElement("div");
					newsel.className = "selector";
					newserver.appendChild(newsel);
					var newicon = document.createElement("p");
					newicon.innerHTML = "` + shortname + `";
					newserver.appendChild(newicon)
					var newtooltip = document.createElement("div");
					newtooltip.className = "tooltip";
					newtooltip.innerHTML = "` + guild.Name + `";
					newserver.appendChild(newtooltip);
					document.getElementById("sidenav").appendChild(newserver);
				`)
			} else {
				wv.Eval(`
					var newserver = document.createElement("div");
					newserver.className = "server";
					newserver.id = "` + guild.ID + `";
					var newsel = document.createElement("div");
					newsel.className = "selector";
					newserver.appendChild(newsel);
					var newicon = document.createElement("img");
					newicon.src = "` + guild.IconURL() + `";
					newserver.appendChild(newicon)
					var newtooltip = document.createElement("div");
					newtooltip.className = "tooltip";
					newtooltip.innerHTML = "` + guild.Name + `";
					newserver.appendChild(newtooltip);
					document.getElementById("sidenav").appendChild(newserver);
				`)
			}
		}
	})
}

func ieUpdate() {
	browser.OpenURL("https://www.microsoft.com/en-us/download/internet-explorer.aspx")
}
