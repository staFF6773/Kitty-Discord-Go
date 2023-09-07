package internal

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler struct {
	cmdMap map[string]*Command
	sync.Mutex
}

func (h *CommandHandler) RegisterCommand(cmd *Command) {
	if _, ok := h.cmdMap[cmd.Name]; ok {
		panic(fmt.Errorf("command exists"))
	}

	h.Lock()
	h.cmdMap[cmd.Name] = cmd
	h.Unlock()
}

// Todo: middlewares
func (h *CommandHandler) Handle(i *discordgo.InteractionCreate) {
	cmd, ok := h.cmdMap[i.ApplicationCommandData().Name]
	if !ok {
		return
	}

	go cmd.Exec(i.Interaction)
}

/*
func (bot *Suplex) RegisterCommands() {
	time.Sleep(time.Second * 3)

	for _, cmd := range bot.Handler.cmdMap {
		fmt.Printf("%#v\n", cmd.AppCmd)

		command := &discordgo.ApplicationCommand{
			Name:        "basic-command",
			Description: "Basic command",
		}

		if _, err := bot.Session.ApplicationCommandCreate(
			bot.Session.State.User.ID,
			"",
			command,
		); err != nil {
			bot.Logger.Error("(*Suplex).RegisterCommands()", cmd.AppCmd.Name, "couldnt register Command", err.Error())
			continue
		}
		fmt.Printf("OK: %s\n", cmd.AppCmd.Name)
	}
}
*/
