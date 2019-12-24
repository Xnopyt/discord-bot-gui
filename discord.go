package main

import (
	"html"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type binder struct{}

type mainBind struct{}

var token string

var ses *discordgo.Session

var currentServer = "HOME"

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

func (m *mainBind) SelectTargetServer(id string) {
	guild, err := ses.Guild(id)
	if err != nil {
		log.Println(err)
		return
	}
	wv.Dispatch(func() {
		wv.Eval(`
		document.getElementsByClassName("server selected")[0].classList.remove("selected");
		document.getElementById("` + id + `").classList.add("selected");
		document.getElementById("servername").innerHTML = "` + html.EscapeString(guild.Name) + `";
		var chancon = document.getElementById("chancontainer");
		chancon.innerHTML = "";
		var head = document.createElement("p");
		head.className = "chanhead";
		head.innerHTML = "TEXT CHANNELS";
		chancon.appendChild(head);
		`)
		chans, _ := ses.GuildChannels(id)
		for _, v := range chans {
			if v.Type == 0 {
				wv.Eval(`
				var chancon = document.getElementById("chancontainer");
				var div = document.createElement("div");
				div.className = "chan";
				var icon = document.createElement("i");
				icon.className = "fas fa-hashtag";
				div.appendChild(icon);
				var para = document.createElement("p");
				para.className = "channame";
				para.innerHTML = "` + html.EscapeString(v.Name) + `";
				div.appendChild(para);
				div.id = "` + v.ID + `";
				chancon.appendChild(div);
				`)
			}
		}
	})
}
