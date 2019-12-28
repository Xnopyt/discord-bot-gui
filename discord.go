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
	ses.AddHandler(recvMsg)
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

func recvMsg(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != currentChannel {
		return
	}
	if m.Type == 7 {
		return
	}
	wv.Dispatch(func() {
		wv.Eval(`createmessage("` + m.ID + `")`)
	})
	wg := &sync.WaitGroup{}
	wg.Add(1)
	processChannelMessage(m, nil, wg)
}

func (m *mainBind) SelectTargetServer(id string) {
	guild, err := ses.Guild(id)
	if err != nil {
		log.Println(err)
		return
	}
	wv.Dispatch(func() {
		wv.Eval(`selectserver("` + id + `", "` + html.EscapeString(guild.Name) + `");`)
		chans, _ := ses.GuildChannels(id)
		var nchan *discordgo.Channel
		i := false
		for _, v := range chans {
			if v.Type == 0 {
				if !i {
					nchan = v
					i = true
				}
				wv.Eval(`addchannel("` + v.ID + `", "` + html.EscapeString(v.Name) + `");`)
			}
		}
		currentServer = id
		m.SetActiveChannel(nchan.ID)
	})
}

func parseTime(m *discordgo.MessageCreate) string {
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
		wv.Eval(`selectchannel("` + id + `", "` + html.EscapeString(channel.Name) + `");
		document.getElementById("mainbox").style.visibility = "hidden";`)
		msgs, err := ses.ChannelMessages(id, 18, "", "", "")
		if err != nil {
			log.Println(err)
			return
		}
		for i := len(msgs)/2 - 1; i >= 0; i-- {
			opp := len(msgs) - 1 - i
			msgs[i], msgs[opp] = msgs[opp], msgs[i]
		}
		wg := &sync.WaitGroup{}
		memberCache, err := ses.GuildMembers(currentServer, "", 1000)
		for _, v := range msgs {
			if v.Type == 7 {
				continue
			}
			wv.Eval(`createmessage("` + v.ID + `")`)
			wg.Add(1)
			go processChannelMessage(&discordgo.MessageCreate{Message: v}, memberCache, wg)
		}
		wg.Wait()
		wv.Eval(`document.getElementById("mainbox").style.visibility = "visible";`)
		currentChannel = id
	})
}

func processChannelMessage(m *discordgo.MessageCreate, cache []*discordgo.Member, wg *sync.WaitGroup) {
	defer wg.Done()
	var uname string
	var member *discordgo.Member
	var err error
	if cache != nil {
		for _, v := range cache {
			if v.User.ID == m.Author.ID {
				member = v
				break
			}
		}
	}
	if member != nil && currentServer != "HOME" {
		member, err = ses.GuildMember(currentServer, m.Author.ID)
	}
	if err == nil && member != nil && currentServer != "HOME" {
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
		wv.Eval(`fillmessage("` + m.ID + `", "` + html.EscapeString(uname) + `", "` + m.Author.AvatarURL("128") + `", "` + parseTime(m) + `", "` + body + `");`)
	})
}

func (m *mainBind) SendMessage(msg string) {
	if currentChannel == "" {
		return
	}
	_, err := ses.ChannelMessageSend(currentChannel, msg)
	if err != nil {
		log.Println(err)
	}
}

func (m *mainBind) LoadDMChannel(id string) {
	channel, err := ses.UserChannelCreate(id)
	if err != nil {
		log.Println(err)
		return
	}
	user, err := ses.User(id)
	if err != nil {
		log.Println(err)
		return
	}
	wv.Dispatch(func() {
		wv.Eval(`selectdmchannel("` + id + `", "` + html.EscapeString(user.Username) + `");
		document.getElementById("mainbox").style.visibility = "hidden";`)
		msgs, err := ses.ChannelMessages(channel.ID, 18, "", "", "")
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
			if v.Type == 7 {
				continue
			}
			wv.Eval(`createmessage("` + v.ID + `")`)
			wg.Add(1)
			go processChannelMessage(&discordgo.MessageCreate{Message: v}, nil, wg)
		}
		wg.Wait()
		wv.Eval(`document.getElementById("mainbox").style.visibility = "visible";`)
		currentChannel = channel.ID
	})
}
