// +build darwin

package main

import (
	"html/template"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/sqweek/dialog"
)

const maxUpload = 8388119

func sendFile(s string) {
	s, err := dialog.File().Load()
	if err != nil {
		return
	}
	f, err := os.Open(s)
	if err != nil {
		wv.Dispatch(func() {
			wv.Eval(`createAlert("Upload Failed", "Failed to open selected file: ` + template.JSEscapeString(err.Error()) + `" );`)
		})
		return
	}
	finfo, _ := f.Stat()
	size := finfo.Size()
	if size > maxUpload {
		wv.Dispatch(func() {
			wv.Eval(`createAlert("Upload Failed", "The selected file exceeds the maximum upload size (8mb).");`)
		})
		return
	}
	var msg discordgo.MessageSend
	msg.Content = ""
	msg.Files = append(msg.Files, &discordgo.File{
		Name:        finfo.Name(),
		ContentType: "application/octet-stream",
		Reader:      f,
	})
	_, err = ses.ChannelMessageSendComplex(currentChannel, &msg)
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
