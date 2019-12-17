package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/zserge/webview"
)

type tokenBind struct{}

var token string

var ses *discordgo.Session

func (t *tokenBind) Store(s string) {
	token = s
}

func login(wv webview.WebView) {
	ses, err := discordgo.New("Bot " + token)
	if err != nil {
		//DO SOMETHING
		return
	}
	err = ses.Open()
	if err != nil {
		//DO SOMETHING
		return
	}
	fmt.Println(ses.State.User.String())
}
