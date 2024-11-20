package handlers

import (
	"discord-bot/internal/handlers/commands"
	image "discord-bot/internal/handlers/commands/images"
	"discord-bot/internal/handlers/commands/matute"
	"discord-bot/internal/store"
	"log"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler struct {
	commands map[string]commands.Command
}

func NewCommandHandler(store *store.ImageStore) *CommandHandler {
	h := &CommandHandler{
		commands: make(map[string]commands.Command),
	}

	// Registrar comandos
	h.registerCommand(image.New(store))
	h.registerCommand(matute.New())

	return h
}

func (h *CommandHandler) registerCommand(cmd commands.Command) {
	h.commands[cmd.Name()] = cmd
	log.Printf("Comando registrado: %s", cmd.Name())
}

func (h *CommandHandler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	commandName := i.ApplicationCommandData().Name
	if command, exists := h.commands[commandName]; exists {
		command.Handle(s, i)
	}
}

func (h *CommandHandler) RegisterCommands(s *discordgo.Session) {
	for _, cmd := range h.commands {
		command := &discordgo.ApplicationCommand{
			Name:        cmd.Name(),
			Description: cmd.Description(),
		}

		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", command)
		if err != nil {
			log.Printf("Error registrando comando %s: %v", cmd.Name(), err)
		}
	}
}
