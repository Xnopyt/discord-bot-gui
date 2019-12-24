package main

import (
	"html"
	"html/template"
	"os"
	"strings"

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

func loadServers() {
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
			wv.Eval(`loadservers("` + html.EscapeString(guild.Name) + `", "` + guild.ID + `", false, "` + html.EscapeString(shortname) + `")`)
		} else {
			wv.Eval(`loadservers("` + html.EscapeString(guild.Name) + `", "` + guild.ID + `", true, "` + guild.IconURL() + `")`)
		}
	}
}

func loadDMMembers() {
	guilds, err := ses.UserGuilds(100, "", "")
	if err != nil {
		panic(err)
	}
	for _, v := range guilds {
		m, err := ses.GuildMembers(v.ID, "", 1000)
		if err == nil {
			for _, x := range m {
				if !x.User.Bot {
					wv.Eval(`loaddmusers("` + html.EscapeString(x.User.Username) + `","` + x.User.ID + `","` + x.User.AvatarURL("128") + `")`)
				}
			}
		}
	}
}
