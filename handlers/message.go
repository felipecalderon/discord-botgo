// handlers/message.go
package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// Estructura para almacenar las URLs de las imágenes
type ImageStore struct {
	Images []string
	mu     sync.RWMutex
}

// Creamos una instancia global
var Store = &ImageStore{
	Images: make([]string, 0),
}

// Método para agregar una imagen
func (s *ImageStore) AddImage(url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Images = append(s.Images, url)
}

// Método para obtener una imagen aleatoria
func (s *ImageStore) GetRandomImage() (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.Images) == 0 {
		return "", fmt.Errorf("no hay imágenes en el canal, aweonao")
	}

	randomIndex := rand.Intn(len(s.Images))
	return s.Images[randomIndex], nil
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignora mensajes del propio bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Monitorear canal específico para guardar imágenes
	monitorChannelID := os.Getenv("MONITOR_CHANNEL_ID")
	if m.ChannelID == monitorChannelID {
		// Guardar imágenes del canal monitoreado
		for _, attachment := range m.Attachments {
			if isImage(attachment) {
				Store.AddImage(attachment.URL)
				log.Printf("Nueva wea de imagen: %s", attachment.URL)
			}
		}
		return
	}

	// Procesar comando .imagen en cualquier canal
	if strings.ToLower(m.Content) == ".imagen" {
		handleImageCommand(s, m)
	}
}

func isImage(attachment *discordgo.MessageAttachment) bool {
	return attachment.Width > 0
}

func handleImageCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	imageURL, err := Store.GetRandomImage()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Esta wea no tiene imagenes")
		return
	}

	// Crear un embed para mostrar la imagen
	embed := &discordgo.MessageEmbed{
		Title: "Toma tu wea random",
		Image: &discordgo.MessageEmbedImage{
			URL: imageURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Esta wea la pidió %s", m.Author.Username),
		},
		Color: 0x00ff00, // Color verde
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		log.Printf("Error enviando la wea de imagen: %v", err)
		s.ChannelMessageSend(m.ChannelID, "Error al enviar la imagen weona")
	}
}
