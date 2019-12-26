package main

import (
	"html"
	"html/template"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type binder struct{}

type mainBind struct{}

var token string

var ses *discordgo.Session

var currentServer = "HOME"
var currentChannel = ""

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
		var nchan *discordgo.Channel
		i := false
		for _, v := range chans {
			if v.Type == 0 {
				if !i {
					nchan = v
					i = true
				}
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
				div.setAttribute("onclick", "bind.setActiveChannel('` + v.ID + `')")
				chancon.appendChild(div);
				`)
			}
		}
		currentServer = id
		m.SetActiveChannel(nchan.ID)
	})
}

func parseTime(m *discordgo.Message) string {
	var ctime string
	times, err := m.Timestamp.Parse()
	if err != nil {
		ctime = "00:00"
	} else {
		times = times.Local()
		hr, mi, _ := times.Clock()
		var min string
		if mi < 10 {
			min = strconv.Itoa(mi)
			min = "0" + min
		} else {
			min = strconv.Itoa(mi)
		}
		ctime = strconv.Itoa(hr) + ":" + min
		y, m, d := times.Date()
		cy, cm, cd := time.Now().Date()
		im := int(m)
		icm := int(cm)
		if y != cy || im != icm || d != cd {
			ctime = strconv.Itoa(d) + "/" + strconv.Itoa(im) + "/" + strconv.Itoa(y)[2:]
		}
	}
	return ctime
}

func (m *mainBind) SetActiveChannel(id string) {
	channel, err := ses.Channel(id)
	if err != nil {
		log.Println(err)
		return
	}
	wv.Dispatch(func() {
		wv.Eval(`
		document.getElementById("infoicon").style.visibility = "visible";
		var title = document.getElementById("channeltitle");
		title.innerHTML = "` + html.EscapeString(channel.Name) + `";
		title.style.visibility = "visible";
		document.getElementById("messageinput").placeholder = "Message #` + html.EscapeString(channel.Name) + `";
		if (document.getElementsByClassName("chan selected")[0]) {
			document.getElementsByClassName("chan selected")[0].classList.remove("selected");
		}
		document.getElementById("` + id + `").classList.add("selected");
		var messages = document.getElementById("messages");
		messages.innerHTML = "";
		var spacer = document.createElement("div");
		spacer.className = "spacer";
		messages.appendChild(spacer);
		`)
		msgs, err := ses.ChannelMessages(id, 30, "", "", "")
		if err != nil {
			log.Println(err)
			return
		}
		for i := len(msgs)/2 - 1; i >= 0; i-- {
			opp := len(msgs) - 1 - i
			msgs[i], msgs[opp] = msgs[opp], msgs[i]
		}
		wg := &sync.WaitGroup{}
		for _, v := range msgs {
			wv.Eval(`
			var messages = document.getElementById("messages");
			var msg = document.createElement("div");
			msg.id = "` + v.ID + `";
			messages.appendChild(msg);
			`)
			wg.Add(1)
			go processChannelMessage(v, wg)
		}
		wg.Wait()
		wv.Eval(`
		document.getElementById("mainbox").style.visibility = "visible";
		`)
		currentChannel = id
	})
}

func processChannelMessage(m *discordgo.Message, wg *sync.WaitGroup) {
	defer wg.Done()
	var uname string
	member, err := ses.GuildMember(m.GuildID, m.Author.ID)
	if err == nil {
		if member.Nick != "" {
			uname = member.Nick
		} else {
			uname = m.Author.Username
		}
	} else {
		uname = m.Author.Username
	}
	for _, z := range m.Attachments {
		if m.Content == "" {
			m.Content = z.URL
		} else {
			m.Content += "\n" + z.URL
		}
	}
	for _, z := range m.Embeds {
		if m.Content != "" {
			m.Content += "\n" + "Embed:"
		}
		if z.Title != "" {
			m.Content += "\n" + z.Title
		}
		if z.Description != "" {
			m.Content += "\n" + z.Description
		}
		if z.URL != "" {
			m.Content += "\n" + z.URL
		}
		if z.Description != "" {
			m.Content += "\n" + z.Description
		}
		if z.Image != nil {
			m.Content += "\n" + z.Image.URL
		}
		if z.Thumbnail != nil {
			m.Content += "\n" + z.Thumbnail.URL
		}
		if z.Video != nil {
			m.Content += "\n" + z.Video.URL
		}
		for _, f := range z.Fields {
			m.Content += "\n" + f.Name + ": " + f.Value
		}
		if z.Provider != nil {
			m.Content += "\n" + "Provider: " + z.Provider.Name + " (" + z.Provider.URL + ")"
		}
		if z.Footer != nil {
			m.Content += "\n" + z.Footer.Text + " " + z.Footer.IconURL
		}
	}
	body, err := m.ContentWithMoreMentionsReplaced(ses)
	if err != nil {
		body = m.ContentWithMentionsReplaced()
	}
	body = html.EscapeString(body)
	body = strings.ReplaceAll(body, "\n", "<br />")
	wv.Dispatch(func() {
		wv.Eval(`
		var msg = document.getElementById("` + m.ID + `");
		msg.className = "message";
		var head = document.createElement("div");
		head.className = "nowrap";
		var ava = document.createElement("img");
		ava.src = "` + m.Author.AvatarURL("128") + `";
		ava.className = "msgavatar";
		debugger;
		head.appendChild(ava);
		var uname = document.createElement("p");
		uname.className = "msguser";
		uname.innerHTML = "` + html.EscapeString(uname) + `";
		head.appendChild(uname);
		var time = document.createElement("p");
		time.className = "msgtime";
		time.innerHTML = "` + parseTime(m) + `";
		head.appendChild(time);
		msg.appendChild(head);
		var body = document.createElement("div");
		body.className = "msgbody";
		body.innerHTML = "` + body + `";
		msg.appendChild(body);
		`)
	})
}
