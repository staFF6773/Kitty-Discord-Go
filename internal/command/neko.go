package command

import (
	"io"
	"net/http"
	"suplex/internal"

	"github.com/bwmarrin/discordgo"
)

const (
	NekoURL = "http://saturn.0x4ba.de:50666/neko"
)

type NekoCommand struct {
	*internal.Suplex
}

func NewNekoCommand(self *internal.Suplex) *internal.Command {
	var (
		cmd = &NekoCommand{self}
	)

	return &internal.Command{
		Exec: cmd.Exec,
		ApplicationCommand: &discordgo.ApplicationCommand{
			Name:        "neko",
			Description: "Nekos UwU",
		},
	}
}

func (c *NekoCommand) Exec(ctx *discordgo.Interaction) {
	//channel is NSFW?

	nekoPicture, err := c.GetNekoPicture()
	if err != nil {
		c.Logger.Error("(*NekoCommand).Exec()", "(*NekoCommand).GetNekoPicture()", err.Error())
		return
	}
	defer nekoPicture.Close()

	if err := c.Session.InteractionRespond(
		ctx,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Title: "Neko",
				Files: []*discordgo.File{
					{
						Name:        "image.png",
						ContentType: "image/png",
						Reader:      nekoPicture,
					},
				},
			},
		},
	); err != nil {
		c.Logger.Error("(*NekoCommand).Exec()", err.Error())
	}
}

func (c *NekoCommand) GetNekoPicture() (io.ReadCloser, error) {
	resp, err := http.Get(NekoURL)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
