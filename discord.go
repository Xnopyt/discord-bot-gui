package main

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
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

type chanCat struct {
	Category *discordgo.Channel   `json:"category"`
	Channels []*discordgo.Channel `json:"channels"`
}

var typing bool

var token string
var proccessingMsg = false

var ses *discordgo.Session

var handlers = [...]interface{}{recvMsg, updateMsg, delMsg, typingStart}

var currentServer = "HOME"
var currentChannel = ""

func connect(s string) {
	token = s
	var err error
	ready := make(chan bool)
	env, ok := os.LookupEnv("DBG_DEBUG_SHARDS")
	if ok {
		shards, err := strconv.Atoi(env)
		if err != nil {
			connectShards(-1)
		} else {
			connectShards(shards)
		}
		goto connected
	}
	ses, err = discordgo.New("Bot " + token)
	if err != nil {
		wv.Dispatch(func() {
			wv.Eval("fail()")
			wv.Eval(`createAlert("` + "Error Creating Session" + `", "` + err.Error() + `");`)
		})
		return
	}
	ses.AddHandler(func(s *discordgo.Session, e *discordgo.Ready) { ready <- true })
	for _, v := range handlers {
		ses.AddHandler(v)
	}
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
connected:
	wv.Dispatch(func() {
		wv.Eval(`document.documentElement.innerHTML="` + template.JSEscapeString(string(MustAsset("ui/main.html"))) + `"`)
	})
	mainSetup()
}

func logout() {
	if shardMan != nil {
		shardMan.stop()
	} else {
		ses.Close()
	}
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
		} else if err.Error() == `HTTP 403 Forbidden, {"message": "Missing Access", "code": 50001}` {
			wv.Dispatch(func() {
				wv.Eval(`createAlert("Failed to Get Guild Members", "Failed to get a list of guild members, please make sure you have Privileged Intents enabled in your bot's settings.")`)
			})
		} else {
			wv.Dispatch(func() {
				wv.Eval(`createAlert("Failed to Get Guild Members", "` + err.Error() + `")`)
			})
		}
	}
	wv.Dispatch(func() {
		wv.Eval(evalQueue)
		wv.Eval(`document.getElementById("blocker").style.display = "none"`)
	})
}

func loadChannels(id string) *discordgo.Channel {
	chans, _ := ses.GuildChannels(id)
	var catChans []*discordgo.Channel
	var categories []chanCat
	var nocat []*discordgo.Channel
	for _, v := range chans {
		switch v.Type {

		case discordgo.ChannelTypeGuildText:
			perms, err := ses.State.UserChannelPermissions(ses.State.User.ID, v.ID)
			if err != nil {
				perms, err = ses.UserChannelPermissions(ses.State.User.ID, v.ID)
				if err != nil {
					continue
				}
			}
			if perms&0x00000400 != 0 {
				if v.ParentID == "" {
					nocat = append(nocat, v)
				} else {
					catChans = append(catChans, v)
				}
			}

		case discordgo.ChannelTypeGuildCategory:
			categories = append(categories, chanCat{Category: v})
		}
	}

	for _, v := range catChans {
		for i, w := range categories {
			if v.ParentID == w.Category.ID {
				categories[i].Channels = append(categories[i].Channels, v)
				break
			}
		}
	}
	var x []chanCat
	for _, v := range categories {
		if len(v.Channels) == 0 {
			continue
		}
		sort.SliceStable(v.Channels, func(i, j int) bool {
			return v.Channels[i].Position < v.Channels[j].Position
		})
		x = append(x, v)
	}
	categories = x
	sort.SliceStable(categories, func(i, j int) bool {
		return categories[i].Category.Position < categories[j].Category.Position
	})
	sort.SliceStable(nocat, func(i, j int) bool {
		return nocat[i].Position < nocat[j].Position
	})

	var cat string
	y, err := json.Marshal(categories)
	if err != nil {
		cat = "null"
	} else {
		cat = string(y)
	}

	var noCat string
	y, err = json.Marshal(nocat)
	if err != nil {
		noCat = "null"
	} else {
		noCat = string(y)
	}

	wv.Dispatch(func() {
		wv.Eval(fmt.Sprintf("addchannels(%q, %q);", noCat, cat))
	})

	if len(nocat) > 0 {
		return nocat[0]
	}
	return categories[0].Channels[0]
}

func selectServer(id string) {
	wv.Dispatch(func() { wv.Eval(`document.getElementById("blocker").style.display = "block"`) })
	time.Sleep(time.Second)
	guild, err := ses.Guild(id)
	if err != nil {
		log.Println(err)
		return
	}
	wv.Dispatch(func() { wv.Eval(fmt.Sprintf(`selectserver(%q, %q);`, id, html.EscapeString(guild.Name))) })
	nchan := loadChannels(id)
	currentServer = id
	setActiveChannel(nchan.ID)
	wv.Dispatch(func() { wv.Eval(`document.getElementById("blocker").style.display = "none"`) })
}

func parseTime(m *discordgo.Message) string {
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

func loadMember(m *discordgo.Member, roles []*discordgo.Role, id string) string {
	perms, err := ses.State.UserChannelPermissions(m.User.ID, id)
	if err != nil {
		perms, err = ses.UserChannelPermissions(m.User.ID, id)
		if err != nil {
			return ""
		}
	}
	if perms&0x00000400 == 0 {
		return ""
	}
	var uname string
	if m.Nick != "" {
		uname = m.Nick
	} else {
		uname = m.User.Username
	}
	var roleColour int
	var colour string
	var hoist string
	var usrroles []*discordgo.Role
	for _, role := range roles {
		for _, rid := range m.Roles {
			if rid == role.ID {
				if role.Color != 0 && roleColour == 0 {
					roleColour = role.Color
				}
				if hoist == "" && role.Hoist {
					hoist = role.ID
				}
				usrroles = append(usrroles, role)
				break
			}
		}
	}
	if roleColour == 0 {
		colour = "null"
	} else {
		colour = fmt.Sprintf("\"#%06x\"", roleColour)
	}
	x, err := json.Marshal(usrroles)
	var usrrolesjson string
	if err == nil {
		usrrolesjson = fmt.Sprintf("%q", x)
	} else {
		usrrolesjson = "null"
	}
	if hoist == "" {
		hoist = "null"
	} else {
		hoist = fmt.Sprintf("%q", hoist)
	}
	return fmt.Sprintf("addmember(%q, %q, %t, %q, %q, %q, %s, %s, %s);\n", html.EscapeString(uname), m.User.AvatarURL("128"), m.User.Bot, m.User.ID, html.EscapeString(m.User.Username), m.User.Discriminator, colour, hoist, usrrolesjson)
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
	if err != nil {
		if err.Error() == `HTTP 403 Forbidden, {"message": "Missing Access", "code": 50001}` {
			wv.Dispatch(func() {
				wv.Eval(`createAlert("Failed to Get Guild Members", "Failed to get a list of guild members, please make sure you have Privileged Intents enabled in your bot's settings.")`)
			})
		} else {
			wv.Dispatch(func() {
				wv.Eval(`createAlert("Failed to Get Guild Members", "` + err.Error() + `")`)
			})
		}
	}
	roles, _ := ses.GuildRoles(currentServer)
	sort.SliceStable(roles, func(i, j int) bool {
		return roles[i].Position > roles[j].Position
	})
	wv.Dispatch(func() {
		wv.Eval(fmt.Sprintf(`selectchannel(%q, %q);
		document.getElementById("members").innerHTML = "";
		resetmembers();`, id, html.EscapeString(channel.Name)))
	})
	var evalQueue string
	rolejson, _ := json.Marshal(roles)
	evalQueue += fmt.Sprintf("loadhoistedroles(%q);\n", string(rolejson))
	for _, v := range memberCache {
		evalQueue += loadMember(v, roles, id)
	}
	evalQueue += "setmembercount();\n"
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
		switch v.Type {

		case discordgo.MessageTypeDefault:
			processChannelMessage(v, nil)
			break

		case discordgo.MessageTypeGuildMemberJoin:
			processMemberJoinMessage(v, nil)
			break

		case discordgo.MessageTypeChannelPinnedMessage:
			processPinnedMessage(v, nil)

		default:
			break
		}
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
		wv.Eval(fmt.Sprintf(`addmember(%q, %q, %t, %q, %q, %q, null, null, null)`, html.EscapeString(ses.State.User.Username), ses.State.User.AvatarURL("128"), ses.State.User.Bot, ses.State.User.ID, html.EscapeString(ses.State.User.Username), ses.State.User.Discriminator))
		for _, v := range channel.Recipients {
			wv.Eval(fmt.Sprintf(`addmember(%q, %q, %t, %q, %q, %q, null, null, null)`, html.EscapeString(v.Username), v.AvatarURL("128"), v.Bot, v.ID, html.EscapeString(v.Username), v.Discriminator))
		}
		wv.Eval(`setmembercount();`)
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
		switch v.Type {

		case discordgo.MessageTypeDefault:
			processChannelMessage(v, nil)
			break

		case discordgo.MessageTypeChannelPinnedMessage:
			processPinnedMessage(v, nil)

		default:
			break
		}
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
