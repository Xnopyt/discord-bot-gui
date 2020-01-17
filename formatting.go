package main

import (
	"html"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func contentWithMentionsFormatted(m *discordgo.MessageCreate) (content string) {
	content = html.EscapeString(m.Content)
	for _, user := range m.Mentions {
		if user.ID == ses.State.User.ID {
			content = strings.NewReplacer(
				"&lt;@"+user.ID+"&gt;", "<div class='selfmention'>@"+user.Username+"</div>",
				"&lt;@!"+user.ID+"&gt;", "<div class='selfmention'>@"+user.Username+"</div>",
			).Replace(content)
			continue
		}
		content = strings.NewReplacer(
			"&lt;@"+user.ID+"&gt;", "<div class='mention'>@"+user.Username+"</div>",
			"&lt;@!"+user.ID+"&gt;", "<div class='mention'>@"+user.Username+"</div>",
		).Replace(content)
	}
	return
}

func contentWithMoreMentionsFormatted(s *discordgo.Session, m *discordgo.MessageCreate) (content string, err error) {
	var patternChannels = regexp.MustCompile("&lt;#[^>]*&gt;")
	content = html.EscapeString(m.Content)

	if !s.StateEnabled {
		content = contentWithMentionsFormatted(m)
		return
	}

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		content = contentWithMentionsFormatted(m)
		return
	}

	for _, user := range m.Mentions {
		nick := user.Username

		member, err := s.State.Member(channel.GuildID, user.ID)
		if err == nil && member.Nick != "" {
			nick = member.Nick
		}
		if user.ID == ses.State.User.ID {
			content = strings.NewReplacer(
				"&lt;@"+user.ID+"&gt;", "<div class='selfmention'>@"+user.Username+"</div>",
				"&lt;@!"+user.ID+"&gt;", "<div class='selfmention'>@"+nick+"</div>",
			).Replace(content)
			continue
		}
		content = strings.NewReplacer(
			"&lt;@"+user.ID+"&gt;", "<div class='mention'>@"+user.Username+"</div>",
			"&lt;@!"+user.ID+"&gt;", "<div class='mention'>@"+nick+"</div>",
		).Replace(content)
	}
	for _, roleID := range m.MentionRoles {
		role, err := s.State.Role(channel.GuildID, roleID)
		if err != nil || !role.Mentionable {
			continue
		}

		content = strings.Replace(content, "&lt;@&"+role.ID+"&gt;", "@"+role.Name, -1)
	}

	content = patternChannels.ReplaceAllStringFunc(content, func(mention string) string {
		channel, err := s.State.Channel(mention[2 : len(mention)-1])
		if err != nil || channel.Type == discordgo.ChannelTypeGuildVoice {
			return mention
		}

		return "#" + channel.Name
	})
	return
}