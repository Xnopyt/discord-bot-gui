package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type binder struct{}

var token string

var ses *discordgo.Session

func (t *binder) Connect(s string) {
	token = s
	ses, err := discordgo.New("Bot " + token)
	if err != nil {
		wv.Dispatch(func() {
			wv.Eval("fail()")
		})
		return
	}
	err = ses.Open()
	if err != nil {
		wv.Dispatch(func() {
			wv.Eval("fail()")
		})
		return
	}
	fmt.Println(ses.State.User.String())
}
