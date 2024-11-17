package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"discord-bot/handlers"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando la wea de archivo .env")
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("Token weon no encontrado")
	}

	// Verificar que existe el MONITOR_CHANNEL_ID
	if os.Getenv("MONITOR_CHANNEL_ID") == "" {
		log.Fatal("Canal weon no encontrado")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creando sesión en la wea de Discord:", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	dg.AddHandler(handlers.MessageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatal("Error conectando la wea:", err)
	}

	fmt.Println("El Bot está corriendo como weon. Presiona CTRL-C para salir.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}
