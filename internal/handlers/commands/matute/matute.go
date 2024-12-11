package matute

import (
	"math/rand"
	"time"

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
	rand.Seed(time.Now().UnixNano())

	// Arreglo de strings para elegir aleatoriamente
	matutePhrases := []string{
		"¡Amo los pitos...! digo, los puzzles!",
		"Puzzles, mi segunda pasión secreta (primera pasión: los pitos)",
		"Los puzles son mi fortaleza, los pitos son mi debilidad",
		"Un puzzle al día mantiene la cordura a raya, un pito al día nunca es suficiente",
	}
	randomPhrase := matutePhrases[rand.Intn(len(matutePhrases))]

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	embeds := []*discordgo.MessageEmbed{
		{
			Description: "[Haz clic acá para jugar a armar pitos, digo puzzles](https://jigsawpuzzles.io/)",
			Color:       0x00FF00, // Opcional: color del embed
		},
	}
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &randomPhrase,
		Embeds:  &embeds,
	})
}
