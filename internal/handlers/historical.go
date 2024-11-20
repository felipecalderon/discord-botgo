package handlers

import (
	"discord-bot/internal/store"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func LoadHistoricalImages(session *discordgo.Session, imageStore *store.ImageStore, channelID string, limit int) error {
	messages, err := session.ChannelMessages(channelID, limit, "", "", "")
	if err != nil {
		return fmt.Errorf("error obteniendo mensajes históricos: %w", err)
	}

	for _, message := range messages {
		for _, attachment := range message.Attachments {
			if attachment.Width > 0 && attachment.Height > 0 {
				imageStore.AddImage(attachment.URL)
				log.Printf("Imagen histórica añadida: %s", attachment.URL)
			}
		}
	}

	return nil
}
