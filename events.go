package main

import (
	"fmt"
	"html"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var guildJoinMessages = [...]string{
	"MEMBER just joined the server - glhf!",
	"MEMBER just joined. Everyone, look busy!",
	"MEMBER just joined. Can I get a heal?",
	"MEMBER joined your party.",
	"MEMBER joined. You must construct additional pylons.",
	"Ermagherd. MEMBER is here.",
	"Welcome, MEMBER. Stay awhile and listen.",
	"Welcome, MEMBER. We hope you brought pizza.",
	"Welcome MEMBER. Leave your weapons by the door.",
	"A wild MEMBER appeared.",
	"Swoooosh. MEMBER just landed.",
	"Brace yourselves. MEMBER just joined the server.",
	"MEMBER just joined. Hide your bananas.",
	"MEMBER just arrived. Seems OP - please nerf.",
	"MEMBER just slid into the server.",
	"A MEMBER has spawned in the server.",
	"Big MEMBER showed up!",
	"Where’s MEMBER? In the server!",
	"MEMBER hopped into the server. Kangaroo!!",
	"MEMBER just showed up. Hold my beer.",
	"Challenger approaching - MEMBER has appeared!",
	"It's a bird! It's a plane! Nevermind, it's just MEMBER.",
	"It's MEMBER! Praise the sun! \\\\[T]/",
	"Never gonna give MEMBER up. Never gonna let MEMBER down.",
	"Ha! MEMBER has joined! You activated my trap card!",
	"Cheers, love! MEMBER's here!",
	"Hey! Listen! MEMBER has joined!",
	"We've been expecting you MEMBER",
	"It's dangerous to go alone, take MEMBER!",
	"MEMBER has joined the server! It's super effective!",
	"Cheers, love! MEMBER is here!",
	"MEMBER is here, as the prophecy foretold.",
	"MEMBER has arrived. Party's over.",
	"Ready player MEMBER",
	"MEMBER is here to kick butt and chew bubblegum. And MEMBER is all out of gum.",
	"Hello. Is it MEMBER you're looking for?",
	"MEMBER has joined. Stay a while and listen!",
	"Roses are red, violets are blue, MEMBER joined this server with you"}

func init() {
	rand.Seed(time.Now().Unix())
}

func recvMsg(s *discordgo.Session, m *discordgo.MessageCreate) {
	for proccessingMsg {
		time.Sleep(time.Second)
	}
	proccessingMsg = true
	defer func() { proccessingMsg = false }()
	if m.ChannelID != currentChannel {
		return
	}
	switch m.Type {

	case discordgo.MessageTypeDefault:
		processChannelMessage(m.Message, nil)
		break

	case discordgo.MessageTypeGuildMemberJoin:
		processMemberJoinMessage(m.Message, nil)
		break

	case discordgo.MessageTypeChannelPinnedMessage:
		processPinnedMessage(m.Message, nil)

	default:
		return
	}
	wv.Dispatch(func() {
		wv.Eval(`var messages = document.getElementsByClassName("messages")[0].querySelector(".simplebar-content-wrapper");
			messages.scrollTop = messages.scrollHeight;`)
	})
}

func updateMsg(s *discordgo.Session, m *discordgo.MessageUpdate) {
	for proccessingMsg {
		time.Sleep(time.Second)
	}
	proccessingMsg = true
	defer func() { proccessingMsg = false }()
	if m.ChannelID != currentChannel {
		return
	}
	if m.Type != discordgo.MessageTypeDefault {
		return
	}
	processChannelMessage(m.Message, nil)
}

func delMsg(s *discordgo.Session, m *discordgo.MessageDelete) {
	if m.ChannelID != currentChannel {
		return
	}
	wv.Dispatch(func() {
		wv.Eval(`document.getElementById("` + m.ID + `").parentNode.removeChild(document.getElementById("` + m.ID + `"));`)
	})
}

func reprocessOnFail(id string, parserFunc func(*discordgo.Message, []*discordgo.Member)) {
	if r := recover(); r != nil {
		msg, err := ses.ChannelMessage(currentChannel, id)
		if err != nil {
			return
		}
		parserFunc(msg, nil)
		wv.Dispatch(func() {
			wv.Eval(`var messages = document.getElementsByClassName("messages")[0].querySelector(".simplebar-content-wrapper");
			messages.scrollTop = messages.scrollHeight;`)
		})
	}
}

func getMember(m []*discordgo.Member, id string) *discordgo.Member {
	if currentServer == "HOME" {
		return nil
	}
	if m != nil {
		for _, v := range m {
			if v.User.ID == id {
				return v
			}
		}
	}
	member, err := ses.GuildMember(currentServer, id)
	if err != nil {
		return member
	}
	return nil
}

func getNick(u *discordgo.User, m *discordgo.Member) string {
	if m == nil {
		if u == nil {
			return ""
		}
		return u.Username
	}
	if m.Nick == "" {
		return u.Username
	}
	return m.Nick
}

func processAttachment(z *discordgo.MessageAttachment, m *discordgo.Message) {
	var isImg = false
	for _, v := range imgMime {
		if strings.HasSuffix(z.URL, v) {
			wv.Dispatch(func() { wv.Eval(fmt.Sprintf(`appendimgattachment(%q, %q);`, m.ID, z.URL)) })
			isImg = true
			break
		}
	}
	if isImg {
		return
	}
	wv.Dispatch(func() { wv.Eval(fmt.Sprintf(`appendattachment(%q, %q, %q);`, m.ID, z.Filename, z.URL)) })
}

func processChannelMessage(m *discordgo.Message, cache []*discordgo.Member) {
	defer reprocessOnFail(m.ID, processChannelMessage)
	member := getMember(cache, m.Author.ID)
	uname := getNick(m.Author, member)
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
		processAttachment(z, m)
	}
}

func processMemberJoinMessage(m *discordgo.Message, cache []*discordgo.Member) {
	defer reprocessOnFail(m.ID, processMemberJoinMessage)
	member := getMember(cache, m.Author.ID)
	uname := getNick(m.Author, member)
	var joinMessage = guildJoinMessages[rand.Intn(len(guildJoinMessages))]
	wv.Dispatch(func() {
		wv.Eval(fmt.Sprintf("createjoinmessage(%q, %q, %q, %q, %q, %q, %q);", m.ID, html.EscapeString(uname), joinMessage, m.Author.ID, m.Author.Discriminator, html.EscapeString(m.Author.Username), parseTime(m)))
	})
}

func processPinnedMessage(m *discordgo.Message, cache []*discordgo.Member) {
	defer reprocessOnFail(m.ID, processPinnedMessage)
	member := getMember(cache, m.Author.ID)
	uname := getNick(m.Author, member)
	wv.Dispatch(func() {
		wv.Eval(fmt.Sprintf("createmessagepinmessage(%q, %q, %q, %q, %q, %q);", m.ID, html.EscapeString(uname), m.Author.ID, m.Author.Discriminator, html.EscapeString(m.Author.Username), parseTime(m)))
	})
}
