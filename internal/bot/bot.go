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
	session    *discordgo.Session
	store      *store.ImageStore
	config     *config.Config
	cmdHandler *handlers.CommandHandler
}

func New(token string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	imageStore := store.NewImageStore(rng)
	cmdHandler := handlers.NewCommandHandler(imageStore)

	return &Bot{
		session:    session,
		store:      imageStore,
		cmdHandler: cmdHandler,
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
	b.session.AddHandler(b.cmdHandler.Handle)
}

func (b *Bot) handleReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Bot está listo y corriendo como weon! Conectado como: %v#%v",
		s.State.User.Username,
		s.State.User.Discriminator)

	// Registrar comandos usando el CommandHandler
	b.cmdHandler.RegisterCommands(s)
}

func (b *Bot) Shutdown() {
	if err := b.session.Close(); err != nil {
		log.Printf("Error cerrando sesión: %v", err)
	}
}
