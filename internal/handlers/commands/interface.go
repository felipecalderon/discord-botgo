package commands

import (
	"github.com/bwmarrin/discordgo"
)

// Command define la interfaz que todos los comandos deben implementar
type Command interface {
	Name() string
	Description() string
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate)
}
