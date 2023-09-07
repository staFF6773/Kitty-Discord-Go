package internal

import "github.com/bwmarrin/discordgo"

type Command struct {
	//? middleware perms
	Exec func(*discordgo.Interaction)
	*discordgo.ApplicationCommand
}
