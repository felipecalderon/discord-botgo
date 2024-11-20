package handlers

import (
	"discord-bot/internal/store"
	"log"

	"github.com/bwmarrin/discordgo"
)

type MessageHandler struct {
	store *store.ImageStore
}

func NewMessageHandler(store *store.ImageStore) *MessageHandler {
	return &MessageHandler{store: store}
}

func (h *MessageHandler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	for _, attachment := range m.Attachments {
		if attachment.Width > 0 && attachment.Height > 0 {
			h.store.AddImage(attachment.URL)
			log.Printf("Imagen añadida: %s", attachment.URL)
		}
	}
}
