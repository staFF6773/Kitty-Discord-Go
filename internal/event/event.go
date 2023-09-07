package event

import (
	"fmt"
	"suplex/internal"

	"github.com/bwmarrin/discordgo"
)

type DefaultHandler struct {
	*internal.Suplex
}

func NewDefaultHandler(self *internal.Suplex) *DefaultHandler {
	return &DefaultHandler{self}
}

func (h *DefaultHandler) Exec(s *discordgo.Session, e *discordgo.Event) {
	h.Logger.Debug(e.Type, fmt.Sprintf("%#v", e.Struct))
}
