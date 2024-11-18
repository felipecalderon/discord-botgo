package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"discord-bot/handlers"

	"github.com/bwmarrin/discordgo"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	rand.Seed(time.Now().UnixNano())

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("ERROR: DISCORD_TOKEN no configurado")
	}

	monitorChannelID := os.Getenv("MONITOR_CHANNEL_ID")
	if monitorChannelID == "" {
		log.Fatal("ERROR: MONITOR_CHANNEL_ID no configurado")
	}

	// Crear nueva sesión de Discord
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creando sesión de Discord: %v", err)
	}

	// Configurar intents necesarios
	dg.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsMessageContent |
		discordgo.IntentsGuilds

	// Registrar handlers
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Bot está listo! Conectado como: %v#%v",
			s.State.User.Username,
			s.State.User.Discriminator)

		// Establecer estado del bot
		err := s.UpdateGameStatus(0, "Usa /imagen para ver fotos random")
		if err != nil {
			log.Printf("Error actualizando estado: %v", err)
		}
	})

	// Registrar comando slash
	cmd := &discordgo.ApplicationCommand{
		Name:        "imagen",
		Description: "Muestra una imagen aleatoria del canal monitoreado",
	}

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			if i.ApplicationCommandData().Name == "imagen" {
				handlers.HandleImageCommand(s, i)
			}
		}
	})

	// Abrir conexión
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error abriendo conexión: %v", err)
	}
	defer dg.Close()

	// Registrar comando en todos los servidores donde está el bot
	guilds, err := dg.UserGuilds(100, "", "", false)
	if err != nil {
		log.Printf("Error obteniendo guilds: %v", err)
	}

	for _, g := range guilds {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, g.ID, cmd)
		if err != nil {
			log.Printf("Error registrando comando en guild %s: %v", g.ID, err)
		}
	}

	// Iniciar servidor HTTP para health check
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Bot está funcionando!")
		})

		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Printf("Error en servidor HTTP: %v", err)
		}
	}()

	log.Println("Bot iniciado. Presiona CTRL-C para detener")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
