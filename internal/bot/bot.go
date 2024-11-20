package bot

import (
	"context"
	"discord-bot/internal/config"
	"discord-bot/internal/handlers"
	"discord-bot/internal/store"
	"log"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session *discordgo.Session
	store   *store.ImageStore
	config  *config.Config
}

func New(cfg *config.Config) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, err
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	imageStore := store.NewImageStore(rng)

	return &Bot{
		session: session,
		store:   imageStore,
		config:  cfg,
	}, nil
}

func (b *Bot) Start(ctx context.Context) error {
	b.session.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsMessageContent |
		discordgo.IntentsGuilds

	// Registrar handlers
	b.registerHandlers()

	if err := b.session.Open(); err != nil {
		return err
	}

	// Cargar imágenes históricas
	return handlers.LoadHistoricalImages(b.session, b.store, b.config.MonitorChannelID, 100)
}

func (b *Bot) registerHandlers() {
	b.session.AddHandler(b.handleReady)
	b.session.AddHandler(handlers.NewMessageHandler(b.store).Handle)
	b.session.AddHandler(handlers.NewCommandHandler(b.store).Handle)
}

func (b *Bot) handleReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Bot está listo! Conectado como: %v#%v",
		s.State.User.Username,
		s.State.User.Discriminator)

	// Registrar comandos
	b.registerCommands()
}

func (b *Bot) registerCommands() {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "imagen",
			Description: "Muestra una imagen aleatoria del canal monitoreado",
		},
		{
			Name:        "matute",
			Description: "GeneralMatute un fan de los puzzles de pitos",
		},
	}

	for _, cmd := range commands {
		_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", cmd)
		if err != nil {
			log.Printf("Error registrando comando %s: %v", cmd.Name, err)
		}
	}
}

func (b *Bot) Shutdown() {
	if err := b.session.Close(); err != nil {
		log.Printf("Error cerrando sesión: %v", err)
	}
}
