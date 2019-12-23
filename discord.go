package main

import (
	"html/template"
	"os"

	"github.com/bwmarrin/discordgo"
)

type binder struct{}

type mainBind struct{}

var token string

var ses *discordgo.Session

func (t *binder) Connect(s string) {
	token = s
	var err error
	ses, err = discordgo.New("Bot " + token)
	if err != nil {
		wv.Dispatch(func() {
			wv.Eval("fail()")
		})
		return
	}
	ready := make(chan bool)
	ses.AddHandler(func(s *discordgo.Session, e *discordgo.Ready) { ready <- true })
	err = ses.Open()
	if err != nil {
		wv.Dispatch(func() {
			wv.Eval("fail()")
		})
		return
	}
	<-ready
	wv.Dispatch(func() {
		wv.Eval(`document.documentElement.innerHTML="` + template.JSEscapeString(string(MustAsset("ui/main.html"))) + `"`)
		wv.Eval("window.external.invoke('mainSetup')")
	})
}

func (m *mainBind) Logout() {
	ses.Close()
	wv.Dispatch(func() {
		wv.Terminate()
	})
	os.Exit(0)
}
