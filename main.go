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
	// Configurar logging detallado
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Inicializar random
	rand.Seed(time.Now().UnixNano())

	// Obtener variables de entorno sin depender de .env
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Printf("WARNING: DISCORD_TOKEN no configurado")
	}

	monitorChannelID := os.Getenv("MONITOR_CHANNEL_ID")
	if monitorChannelID == "" {
		log.Printf("WARNING: MONITOR_CHANNEL_ID no configurado")
	}

	// Iniciar servidor HTTP primero para el health check
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Bot is running")
		})

		port := os.Getenv("PORT")
		if port == "" {
			port = "8080" // Puerto por defecto si no está configurado
		}

		server := &http.Server{
			Addr:    ":" + port,
			Handler: mux,
		}

		log.Printf("Iniciando servidor HTTP en puerto %s", port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("INFO: Servidor HTTP: %v", err)
		}
	}()

	// Solo iniciar el bot si tenemos las variables necesarias
	if token != "" {
		dg, err := discordgo.New("Bot " + token)
		if err != nil {
			log.Printf("Error creando sesión de Discord: %v", err)
			return
		}

		// Configurar intents
		dg.Identify.Intents = discordgo.IntentsGuildMessages |
			discordgo.IntentsMessageContent |
			discordgo.IntentsGuilds

		// Agregar handler de ready para confirmar conexión
		dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
			log.Printf("Bot está listo! Conectado como: %v#%v",
				s.State.User.Username,
				s.State.User.Discriminator)
		})

		dg.AddHandler(handlers.MessageCreate)

		err = dg.Open()
		if err != nil {
			log.Printf("Error abriendo conexión Discord: %v", err)
			return
		}

		defer dg.Close()

		log.Println("Bot iniciado correctamente")
	} else {
		log.Println("Bot no iniciado por falta de token")
	}

	// Mantener el programa corriendo
	log.Println("Servicio iniciado. Presiona CTRL-C para detener")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Cerrando servicio...")
}
