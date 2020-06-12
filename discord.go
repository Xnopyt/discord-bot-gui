package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var imgMime = []string{
	".bmp",
	".gif",
	".jpe",
	".jpeg",
	".jpg",
	".svg",
	".ico",
	".png",
}

type fileAttachment struct {
	Data string `json:"data"`
	Name string `json:"name"`
	Mime string `json:"mime"`
}

type emojiAliases struct {
	Emojis []struct {
		Aliases []string `json:"aliases"`
		Unicode string   `json:"unicode"`
	} `json:"emojis"`
}

const maxUpload = 8388119

var eAliases emojiAliases

var typing bool

var token string
var proccessingMsg = false

var ses *discordgo.Session

var currentServer = "HOME"
var currentChannel = ""

func connect(s string) {
	token = s
	var err error
	ses, err = discordgo.New("Bot " + token)
	if err != nil {
		wv.Dispatch(func() {
			wv.Eval("fail()")
			wv.Eval(`createAlert("` + "Error Creating Session" + `", "` + err.Error() + `");`)
		})
		return
	}
	ready := make(chan bool)
	ses.AddHandler(func(s *discordgo.Session, e *discordgo.Ready) { ready <- true })
	ses.AddHandler(recvMsg)
	ses.AddHandler(updateMsg)
	ses.AddHandler(delMsg)
	err = ses.Open()
	if err != nil {
		wv.Dispatch(func() {
			wv.Eval("fail()")
			if err.Error() != "websocket: close 4004: Authentication failed." {
				wv.Eval(`createAlert("` + "Error Opening Session" + `", "` + err.Error() + `");`)
			}
		})
		return
	}
	<-ready
	wv.Dispatch(func() {
		wv.Eval(`document.documentElement.innerHTML="` + template.JSEscapeString(string(MustAsset("ui/main.html"))) + `"`)
	})
	mainSetup()
}

func logout() {
	ses.Close()
	wv.Terminate()
	os.Exit(0)
}

func loadServers() {
	guilds, err := ses.UserGuilds(100, "", "")
	if err != nil {
		panic(err)
	}
	var evalQueue string
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
			evalQueue += fmt.Sprintf("loadservers(%q, %q, %t, %q);\n", html.EscapeString(guild.Name), guild.ID, false, html.EscapeString(shortname))
		} else {
			evalQueue += fmt.Sprintf("loadservers(%q, %q, %t, %q)\n", html.EscapeString(guild.Name), guild.ID, true, guild.IconURL())
		}
	}
	wv.Dispatch(func() { wv.Eval(evalQueue) })
}

func loadDMMembers() {
	wv.Dispatch(func() { wv.Eval(`document.getElementById("blocker").style.display = "block"`) })
	time.Sleep(time.Second)
	guilds, err := ses.UserGuilds(100, "", "")
	if err != nil {
		panic(err)
	}
	var evalQueue string
	for _, v := range guilds {
		m, err := ses.GuildMembers(v.ID, "", 1000)
		if err == nil {
			for _, x := range m {
				if !x.User.Bot {
					evalQueue += fmt.Sprintf("loaddmusers(%q,%q,%q);\n", html.EscapeString(x.User.Username), x.User.ID, x.User.AvatarURL("128"))
				}
			}
		}
	}
	wv.Dispatch(func() {
		wv.Eval(evalQueue)
		wv.Eval(`document.getElementById("blocker").style.display = "none"`)
	})
}

func recvMsg(s *discordgo.Session, m *discordgo.MessageCreate) {
	for proccessingMsg {
		time.Sleep(time.Second)
	}
	proccessingMsg = true
	if m.ChannelID != currentChannel {
		return
	}
	if m.Type == 7 {
		return
	}
	processChannelMessage(m, nil)
	wv.Dispatch(func() {
		wv.Eval(`var messages = document.getElementsByClassName("messages")[0].querySelector(".simplebar-content-wrapper");
		messages.scrollTop = messages.scrollHeight;`)
		proccessingMsg = false
	})
}

func updateMsg(s *discordgo.Session, m *discordgo.MessageUpdate) {
	for proccessingMsg {
		time.Sleep(time.Second)
	}
	proccessingMsg = true
	if m.ChannelID != currentChannel {
		return
	}
	if m.Type == 7 {
		return
	}
	processChannelMessage(&discordgo.MessageCreate{Message: m.Message}, nil)
	proccessingMsg = false
}

func delMsg(s *discordgo.Session, m *discordgo.MessageDelete) {
	if m.ChannelID != currentChannel {
		return
	}
	if m.Type == 7 {
		return
	}
	wv.Dispatch(func() {
		wv.Eval(`document.getElementById("` + m.ID + `").parentNode.removeChild(document.getElementById("` + m.ID + `"));`)
	})
}

func selectTargetServer(id string) {
	wv.Dispatch(func() { wv.Eval(`document.getElementById("blocker").style.display = "block"`) })
	time.Sleep(time.Second)
	guild, err := ses.Guild(id)
	if err != nil {
		log.Println(err)
		return
	}
	wv.Dispatch(func() { wv.Eval(fmt.Sprintf(`selectserver(%q, %q);`, id, html.EscapeString(guild.Name))) })
	chans, _ := ses.GuildChannels(id)
	var nchan *discordgo.Channel
	i := false
	var evalQueue string
	for _, v := range chans {
		if v.Type == 0 {
			if !i {
				nchan = v
				i = true
			}
			perms, err := ses.State.UserChannelPermissions(ses.State.User.ID, v.ID)
			if err != nil {
				continue
			}
			if perms&0x00000400 != 0 {
				evalQueue += fmt.Sprintf("addchannel(%q, %q);\n", v.ID, html.EscapeString(v.Name))
			}
		}
	}
	wv.Dispatch(func() {
		wv.Eval(evalQueue)
	})
	currentServer = id
	setActiveChannel(nchan.ID)
	wv.Dispatch(func() { wv.Eval(`document.getElementById("blocker").style.display = "none"`) })
}

func parseTime(m *discordgo.MessageCreate) string {
	var ctime string
	times, err := discordgo.SnowflakeTimestamp(m.ID)
	if err != nil {
		ctime = "Unable to Parse Timestamp"
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
		im := int(m)
		ctime = strconv.Itoa(d) + "/" + strconv.Itoa(im) + "/" + strconv.Itoa(y)[2:] + " at " + ctime
	}
	return ctime
}

func setActiveChannel(id string) {
	wv.Dispatch(func() {
		wv.Eval(`document.getElementById("blocker").style.display = "block";
			document.getElementById("mainbox").style.visibility = "hidden";
			document.getElementById("mainbox").style.display = "inline-block";`)
	})
	time.Sleep(time.Second)
	channel, err := ses.Channel(id)
	if err != nil {
		log.Println(err)
		wv.Dispatch(func() { wv.Eval(`document.getElementById("blocker").style.display = "none"`) })
		return
	}
	memberCache, err := ses.GuildMembers(currentServer, "", 1000)
	roles, _ := ses.GuildRoles(currentServer)
	sort.SliceStable(roles, func(i, j int) bool {
		return roles[i].Position > roles[j].Position
	})
	wv.Dispatch(func() {
		wv.Eval(fmt.Sprintf(`selectchannel(%q, %q);
		document.getElementById("members").innerHTML = "";
		resetmembers();`, id, html.EscapeString(channel.Name)))
	})
	var i = 0
	var evalQueue string
	for _, v := range memberCache {
		perms, err := ses.State.UserChannelPermissions(v.User.ID, id)
		if err != nil {
			continue
		}
		if perms&0x00000400 != 0 {
			i++
			var uname string
			if v.Nick != "" {
				uname = v.Nick
			} else {
				uname = v.User.Username
			}
			var roleColour int
			var colour string
			for _, role := range roles {
				if roleColour != 0 {
					break
				}
				for _, rid := range v.Roles {
					if rid == role.ID && role.Color != 0 {
						roleColour = role.Color
						break
					}
				}
			}
			if roleColour == 0 {
				colour = "null"
			} else {
				colour = fmt.Sprintf("\"#%06x\"", roleColour)
			}
			evalQueue += fmt.Sprintf("addmember(%q, %q, %t, %q, %q, %q, %s);\n", html.EscapeString(uname), v.User.AvatarURL("128"), v.User.Bot, v.User.ID, html.EscapeString(v.User.Username), v.User.Discriminator, colour)
		}
	}
	evalQueue += fmt.Sprintf("setmembercount('%d');\n", i)
	wv.Dispatch(func() {
		wv.Eval(evalQueue)
	})
	msgs, err := ses.ChannelMessages(id, 18, "", "", "")
	if err != nil {
		log.Println(err)
		return
	}
	for i := len(msgs)/2 - 1; i >= 0; i-- {
		opp := len(msgs) - 1 - i
		msgs[i], msgs[opp] = msgs[opp], msgs[i]
	}
	for _, v := range msgs {
		if v.Type == 7 {
			continue
		}
		processChannelMessage(&discordgo.MessageCreate{Message: v}, memberCache)
	}
	time.Sleep(time.Second)
	wv.Dispatch(func() {
		wv.Eval(`var messages = document.getElementsByClassName("messages")[0].querySelector(".simplebar-content-wrapper");
		messages.scrollTop = messages.scrollHeight;
		document.getElementById("mainbox").style.visibility = "visible";
		document.getElementById("blocker").style.display = "none"`)
	})
	currentChannel = id
}

func processChannelMessage(m *discordgo.MessageCreate, cache []*discordgo.Member) {
	defer func(id string) {
		if r := recover(); r != nil {
			msg, err := ses.ChannelMessage(currentChannel, id)
			if err != nil {
				return
			}
			processChannelMessage(&discordgo.MessageCreate{Message: msg}, nil)
			wv.Dispatch(func() {
				wv.Eval(`var messages = document.getElementsByClassName("messages")[0].querySelector(".simplebar-content-wrapper");
				messages.scrollTop = messages.scrollHeight;`)
			})
		}
	}(m.ID)
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
	var embeds string
	for _, z := range m.Embeds {
		embeds += processEmbed(z, m) + `
		document.getElementById("` + m.ID + `").appendChild(div);
		`
	}
	body := parseMarkdownAndMentions(m)
	body = strings.ReplaceAll(body, "\n", "<br />")
	var selfmention = false
	if strings.Contains(body, "<div class='selfmention'") {
		selfmention = true
	}
	wv.Dispatch(func() {
		wv.Eval(fmt.Sprintf(`fillmessage(%q, %q, %q, %q, %q, %t, %t, %q, %q, %q);`, m.ID, html.EscapeString(uname), m.Author.AvatarURL("128"), parseTime(m), url.QueryEscape(body), selfmention, m.Author.Bot, m.Author.ID, m.Author.Discriminator, html.EscapeString(m.Author.Username)))
		wv.Eval(embeds)
	})
	for _, z := range m.Attachments {
		var isImg = false
		for _, v := range imgMime {
			if strings.HasSuffix(z.URL, v) {
				wv.Dispatch(func() { wv.Eval(fmt.Sprintf(`appendimgattachment(%q, %q);`, m.ID, z.URL)) })
				isImg = true
				break
			}
		}
		if isImg {
			continue
		}
		wv.Dispatch(func() { wv.Eval(fmt.Sprintf(`appendattachment(%q, %q, %q);`, m.ID, z.Filename, z.URL)) })
	}
}

func sendMessage(msg string) {
	go func() {
		if currentChannel == "" {
			return
		}
		_, err := ses.ChannelMessageSend(currentChannel, msg)
		if err != nil {
			log.Println(err)
		}
	}()
}

func loadDMChannel(id string) {
	wv.Dispatch(func() {
		wv.Eval(`document.getElementById("blocker").style.display = "block";
				document.getElementById("mainbox").style.visibility = "hidden";
				document.getElementById("mainbox").style.display = "inline-block";`)
	})
	channel, err := ses.UserChannelCreate(id)
	if err != nil {
		log.Println(err)
		wv.Dispatch(func() { wv.Eval(`document.getElementById("blocker").style.display = "none"`) })
		return
	}
	user, err := ses.User(id)
	if err != nil {
		log.Println(err)
		return
	}
	wv.Dispatch(func() {
		wv.Eval(fmt.Sprintf(`selectdmchannel(%q, %q);`, id, html.EscapeString(user.Username)))
		wv.Eval(`resetmembers();`)
		wv.Eval(fmt.Sprintf(`addmember(%q, %q, %t, %q, %q, %q, null)`, html.EscapeString(ses.State.User.Username), ses.State.User.AvatarURL("128"), ses.State.User.Bot, ses.State.User.ID, html.EscapeString(ses.State.User.Username), ses.State.User.Discriminator))
		for _, v := range channel.Recipients {
			wv.Eval(fmt.Sprintf(`addmember(%q, %q, %t, %q, %q, %q, null)`, html.EscapeString(v.Username), v.AvatarURL("128"), v.Bot, v.ID, html.EscapeString(v.Username), v.Discriminator))
		}
		wv.Eval(fmt.Sprintf(`setmembercount("%d");`, len(channel.Recipients)+1))
	})
	msgs, err := ses.ChannelMessages(channel.ID, 18, "", "", "")
	if err != nil {
		log.Println(err)
		return
	}
	for i := len(msgs)/2 - 1; i >= 0; i-- {
		opp := len(msgs) - 1 - i
		msgs[i], msgs[opp] = msgs[opp], msgs[i]
	}
	for _, v := range msgs {
		if v.Type == 7 {
			continue
		}
		processChannelMessage(&discordgo.MessageCreate{Message: v}, nil)
	}
	time.Sleep(time.Second)
	wv.Dispatch(func() {
		wv.Eval(`var messages = document.getElementsByClassName("messages")[0].querySelector(".simplebar-content-wrapper");
	messages.scrollTop = messages.scrollHeight;
	document.getElementById("mainbox").style.display = "inline-block";
	document.getElementById("mainbox").style.visibility = "visible";
	document.getElementById("blocker").style.display = "none"`)
	})
	currentChannel = channel.ID
}

func sendFile(s string) {
	var file fileAttachment
	json.Unmarshal([]byte(s), &file)
	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(file.Data))
	var msg discordgo.MessageSend
	msg.Content = ""
	msg.Files = append(msg.Files, &discordgo.File{
		Name:        file.Name,
		ContentType: file.Mime,
		Reader:      decoder,
	})
	_, err := ses.ChannelMessageSendComplex(currentChannel, &msg)
	if err != nil {
		wv.Dispatch(func() { wv.Eval(`createAlert("Upload Failed", "` + template.JSEscapeString(err.Error()) + `");`) })
	}
}

func updateTyping() {
	go func() {
		if typing {
			return
		}
		ses.ChannelTyping(currentChannel)
		typing = true
		time.Sleep(time.Second * 3)
		typing = false
	}()
}

func deleteMessage(id string) string {
	err := ses.ChannelMessageDelete(currentChannel, id)
	if err != nil {
		return err.Error()
	}
	return ""
}
