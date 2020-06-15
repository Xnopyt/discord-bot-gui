package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/atotto/clipboard"
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

func httpGet(url string) (body []byte) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("could not download " + url)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("could not download " + url)
	}
	return
}

func injectCSS(css []byte) {
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
	})("%s")`, template.JSEscapeString(string(css))))
}

func injectJS(js []byte) {
	wv.Eval(fmt.Sprintf(`
			var script = document.createElement('script');
			var head = document.head || document.getElementsByTagName('head')[0];
			script.src="data:application/javascript;base64,%s"
			head.appendChild(script)`, base64.StdEncoding.EncodeToString(js)))
}

func injectCSSFromURL(url string) {
	injectCSS(httpGet(url))
}

func injectJSFromURL(url string) {
	injectJS(httpGet(url))
}

func loginSetup() {
	wv.Dispatch(func() {
		injectJS(MustAsset("ui/js/login.js"))
		if runtime.GOOS == "darwin" {
			injectJS(MustAsset("ui/js/darwinClipboard.js"))
		}
		injectCSS(MustAsset("ui/login.css"))
		wv.Eval(fmt.Sprintf(`
			document.body.background = "data:image/jpg;base64,%s"`, base64.StdEncoding.EncodeToString(MustAsset("ui/assets/loginbg.jpg"))))
	})
}

func mainSetup() {
	wv.Dispatch(func() {
		injectCSSFromURL("https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@9.17.1/build/styles/androidstudio.min.css")
		injectCSSFromURL("https://cdnjs.cloudflare.com/ajax/libs/simplebar/5.2.0/simplebar.min.css")
		injectCSS(MustAsset("ui/main.css"))
		if runtime.GOOS == "windows" {
			injectCSS([]byte(template.JSEscapeString(`
			.infobar .chantitle {
				transform: none;
			}

			.infobar .fa-hashtag, .infobar .fa-at {
				transform: translateY(-15px);
			}

			.chan .fa-hashtag {
				transform: translateY(-3px);
			}

			.memberbar .memberbot {
				transform: translateY(-16px);
				padding-left: 2px;
			}

			.message .msgbot {
				transform: translateY(-24px);
				padding-left: 2px;
			}

			.emojiselect .fa-grin {
				transform: none;
			}

			.actionbar .dmusername {
				transform: translateY(-5px);
			}

			.fileupload .fa-plus-circle {
				transform: none;
			}

			.attachment p {
				transform: none;
			}
			`)))
		}
		injectJSFromURL("https://cdn.jsdelivr.net/npm/@joeattardi/emoji-button@2.8.2/dist/index.min.js")
		injectJSFromURL("https://twemoji.maxcdn.com/v/latest/twemoji.min.js")
		if runtime.GOOS == "windows" {
			time.Sleep(time.Second)
		}
		injectJS(MustAsset("ui/js/main.js"))
		injectJSFromURL("https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@9.17.1/build/highlight.min.js")
		injectJSFromURL("https://cdnjs.cloudflare.com/ajax/libs/simplebar/5.2.0/simplebar.min.js")
		injectCSS(MustAsset("ui/emoji-picker.css"))
		wv.Eval(fmt.Sprintf(`
			document.getElementById("cname").innerHTML = %q;
			document.getElementById("cdiscriminator").innerHTML = '#%s';
			document.getElementById("cavatar").src = %q;
		`, html.EscapeString(ses.State.User.Username), ses.State.User.Discriminator, ses.State.User.AvatarURL("128")))
		wv.SetTitle(fmt.Sprintf(`Discord Bot GUI - %s#%s`, ses.State.User.Username, ses.State.User.Discriminator))
	})
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

func readClipboard() string {
	clip, err := clipboard.ReadAll()
	if err != nil {
		log.Println("Error Reading Clipboard: " + err.Error())
	}
	return clip
}

func writeClipboard(clip string) {
	err := clipboard.WriteAll(clip)
	if err != nil {
		log.Println("Error Writing Clipboard: " + err.Error())
	}
}
