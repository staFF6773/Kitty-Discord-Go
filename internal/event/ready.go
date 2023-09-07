package event

import (
	"suplex/internal"

	"github.com/bwmarrin/discordgo"
)

type ReadyHandler struct {
	*internal.Suplex
}

func NewReadyHandler(self *internal.Suplex) *ReadyHandler {
	return &ReadyHandler{self}
}

func (h *ReadyHandler) Exec(s *discordgo.Session, e *discordgo.Ready) {
	h.Logger.Info("Ready", "Bot up as", e.User.Username)

	// Register All registered Slash-Commands
	h.ReregisterAllCommands()
}
