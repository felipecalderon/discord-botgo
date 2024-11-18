package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type ImageStore struct {
	Images []string
	rng    *rand.Rand
	mu     sync.RWMutex
}

var Store *ImageStore

// InitializeStore inicializa el store con un generador de números aleatorios
func InitializeStore(rng *rand.Rand) {
	Store = &ImageStore{
		Images: make([]string, 0),
		rng:    rng,
	}
}

func (s *ImageStore) AddImage(url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Images = append(s.Images, url)
	log.Printf("Nueva imagen agregada: %s", url)
}

func (s *ImageStore) GetRandomImage() (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.Images) == 0 {
		return "", fmt.Errorf("no hay imágenes disponibles")
	}

	return s.Images[s.rng.Intn(len(s.Images))], nil
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	monitorChannelID := os.Getenv("MONITOR_CHANNEL_ID")
	if m.ChannelID == monitorChannelID {
		for _, attachment := range m.Attachments {
			if isImage(attachment) {
				Store.AddImage(attachment.URL)
			}
		}
	}
}

func HandleImageCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Función helper para convertir string a *string
	strPtr := func(s string) *string {
		return &s
	}

	// Responder inmediatamente para evitar timeout
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	imageURL, err := Store.GetRandomImage()
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: strPtr("No hay imágenes disponibles en este momento"),
		})
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "¡Aquí tienes tu imagen random!",
		Image: &discordgo.MessageEmbedImage{
			URL: imageURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Solicitada por %s", i.Member.User.Username),
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

func isImage(attachment *discordgo.MessageAttachment) bool {
	return attachment.Width > 0
}
