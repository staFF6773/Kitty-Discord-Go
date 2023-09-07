package event

import (
	"suplex/internal"

	"github.com/bwmarrin/discordgo"
)

type InteractionAdd struct {
	*internal.Suplex
}

func NewInteractionAdd(self *internal.Suplex) *InteractionAdd {
	return &InteractionAdd{self}
}

func (h *InteractionAdd) Exec(s *discordgo.Session, e *discordgo.InteractionCreate) {
	h.Handler.Handle(e)
}
