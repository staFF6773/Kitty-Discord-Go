package command

import (
	"suplex/internal"

	"github.com/bwmarrin/discordgo"
)

type TestCommand struct {
	*internal.Suplex
}

func NewTestCommand(self *internal.Suplex) *internal.Command {
	var cmd = &TestCommand{self}
	return &internal.Command{
		Exec: cmd.Exec,
		ApplicationCommand: &discordgo.ApplicationCommand{
			Name:        "test",
			Description: "test command",
		},
	}
}

func (c *TestCommand) Exec(ctx *discordgo.Interaction) {

	if err := c.Session.InteractionRespond(
		ctx,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "bububu",
			},
		},
	); err != nil {
		c.Logger.Error("(*TestCommand).Exec()", err.Error())
	}
}
