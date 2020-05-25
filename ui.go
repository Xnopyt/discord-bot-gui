package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"runtime"
	"time"
)

type uiMsg struct {
	Type    string `json:"type"`
	Content string `json:"content"`
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

func loginSetup() {
	wv.Dispatch(func() {
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
		wv.Eval(fmt.Sprintf(`
			document.body.background = "data:image/jpg;base64,%s"`, base64.StdEncoding.EncodeToString(MustAsset("ui/assets/loginbg.jpg"))))
	})
}

func mainSetup() {
	wv.Dispatch(func() {
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
		if runtime.GOOS == "windows" {
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
			})("%s")`, template.JSEscapeString(`
			.infobar .chantitle {
				transform: none;
			}

			.infobar .fa-hashtag, .infobar .fa-at {
				transform: translateY(-15px);
			}

			.chan .fa-hashtag {
				transform: translateY(-3px);
			}
			`)))
		}
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
	})
	wv.SetTitle(fmt.Sprintf(`Discord Bot GUI - %s#%s`, html.EscapeString(ses.State.User.Username), ses.State.User.Discriminator))
	time.Sleep(time.Second)
	loadServers()
	loadDMMembers()
}

func home() {
	currentServer = "HOME"
	currentChannel = ""
	wv.Dispatch(func() {
		wv.Eval(`loadhome()`)
	})
	loadDMMembers()
}
