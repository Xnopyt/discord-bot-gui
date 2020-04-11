package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type uiMsg struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

var routes = []route{
	route{"Login", "GET", "/", loginPage}, route{"Lbg", "GET", "/loginbg.jpg", lbg}, route{"DefAva", "GET", "/default.png", defaultavatar},
}

var wvCallbacks map[string]func()

func init() {
	json.Unmarshal(MustAsset("ui/assets/emojialiases.json"), &eAliases)
	wvCallbacks = make(map[string]func())

	wvCallbacks["loginSetup"] = loginSetup
	wvCallbacks["home"] = home
	wvCallbacks["logout"] = logout
	wvCallbacks["updateTyping"] = updateTyping
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

func loginSetup() {
	wv.Eval(fmt.Sprintf(`
		var script = document.createElement('script');
		var head = document.head || document.getElementsByTagName('head')[0];
		script.src="data:application/javascript;base64,%s"
		head.appendChild(script)`, base64.StdEncoding.EncodeToString(MustAsset("ui/js/login.js"))))
	wv.Eval(fmt.Sprintf(`(function(css){
		var style = document.createElement('style');
		var head = document.head || document.getElementsByTagName('head')[0];
		style.setAttribute('type', 'text/css');
		if (style.styleSheet) {
			style.styleSheet.cssText = css;
		} else {
			style.appendChild(document.createTextNode(css));
		}
		head.appendChild(style);
	})("%s")`, template.JSEscapeString(string(MustAsset("ui/login.css")))))
}

func mainSetup() {
	wv.Eval(fmt.Sprintf(`(function(css){
		var style = document.createElement('style');
		var head = document.head || document.getElementsByTagName('head')[0];
		style.setAttribute('type', 'text/css');
		if (style.styleSheet) {
			style.styleSheet.cssText = css;
		} else {
			style.appendChild(document.createTextNode(css));
		}
		head.appendChild(style);
	})("%s")`, template.JSEscapeString(string(MustAsset("ui/main.css")))))
	wv.Eval(fmt.Sprintf(`(function(css){
		var style = document.createElement('style');
		var head = document.head || document.getElementsByTagName('head')[0];
		style.setAttribute('type', 'text/css');
		if (style.styleSheet) {
			style.styleSheet.cssText = css;
		} else {
			style.appendChild(document.createTextNode(css));
		}
		head.appendChild(style);
	})("%s")`, template.JSEscapeString(string(MustAsset("ui/emoji-picker.css")))))
	wv.Eval(fmt.Sprintf(`
		var script = document.createElement('script');
		var head = document.head || document.getElementsByTagName('head')[0];
		script.src="data:application/javascript;base64,%s"
		head.appendChild(script)`, base64.StdEncoding.EncodeToString(MustAsset("ui/js/main.js"))))
	wv.Eval(fmt.Sprintf(`
		document.getElementById("cname").innerHTML = %q;
		document.getElementById("cdiscriminator").innerHTML = '#%s';
		document.getElementById("cavatar").src = %q;
	`, html.EscapeString(ses.State.User.Username), ses.State.User.Discriminator, ses.State.User.AvatarURL("128")))
	wv.Eval(fmt.Sprintf(`document.title = "Discord Bot GUI - %s#%s";`, html.EscapeString(ses.State.User.Username), ses.State.User.Discriminator))
	time.Sleep(time.Second)
	loadServers()
	loadDMMembers()
}

func home() {
	currentServer = "HOME"
	currentChannel = ""
	wv.Eval(`loadhome()`)
	loadDMMembers()
}
