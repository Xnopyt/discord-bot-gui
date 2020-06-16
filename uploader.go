// +build !darwin

package main

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"strings"

	"github.com/bwmarrin/discordgo"
)

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
		wv.Dispatch(func() {
			wv.Eval(`
				var err = "` + template.JSEscapeString(err.Error()) + `";
				var x = err.split(",");
				x.shift();
				x = x.join(",");
				try {
					err = JSON.parse(x).message;
				} catch (e) {}
				createAlert("Upload Failed", err);`)
		})
	}
}
