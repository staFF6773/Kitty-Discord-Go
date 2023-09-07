package event

import (
	"suplex/internal"

	"github.com/bwmarrin/discordgo"
)

type MessageAddHandler struct {
	*internal.Suplex
}

func NewMessageAddHandler(self *internal.Suplex) *MessageAddHandler {
	return &MessageAddHandler{self}
}

func (h *MessageAddHandler) Exec(s *discordgo.Session, e *discordgo.MessageCreate) {}
