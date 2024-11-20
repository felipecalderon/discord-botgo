package matute

import (
	"github.com/bwmarrin/discordgo"
)

type MatuteCommand struct{}

func New() *MatuteCommand {
	return &MatuteCommand{}
}

func (c *MatuteCommand) Name() string {
	return "matute"
}

func (c *MatuteCommand) Description() string {
	return "GeneralMatute se devora los puzzles (y los pitos)"
}

func (c *MatuteCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Implementación del comando matute
	strPtr := func(s string) *string { return &s }

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: strPtr("¡Amo los pitos...! digo, los puzzles!"),
	})
}
