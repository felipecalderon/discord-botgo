package image

import (
	"discord-bot/internal/store"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type ImageCommand struct {
	store *store.ImageStore
}

func New(store *store.ImageStore) *ImageCommand {
	return &ImageCommand{
		store: store,
	}
}

func (c *ImageCommand) Name() string {
	return "imagen"
}

func (c *ImageCommand) Description() string {
	return "Muestra una imagen aleatoria del canal monitoreado"
}

func (c *ImageCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	strPtr := func(s string) *string { return &s }

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	imageURL, err := c.store.GetRandomImage()
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: strPtr("No hay imágenes disponibles en este momento"),
		})
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "¡Toma tu wea ctm!",
		Image: &discordgo.MessageEmbedImage{
			URL: imageURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("El CTM que pidió esta wea fue %s", i.Member.User.Username),
		},
		Color: 0x00ff00,
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})

	if err != nil {
		log.Printf("Error enviando imagen: %v", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: strPtr("Ocurrió un error al enviar la imagen"),
		})
	}
}
