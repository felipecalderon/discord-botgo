package main

import (
	"context"
	"discord-bot/config"
	"discord-bot/internal/bot"
	"discord-bot/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env:", err)
	}
	ctx := context.Background()

	var cfg *config.Config
	var err error
	var port string
	// cargar desde archivo de configuración (local)
	if _, err := os.Stat("config/config.json"); err == nil {
		cfg, err = config.LoadConfig("config/config.json")
		if err != nil {
			log.Fatalf("Error cargando configuración local: %v", err)
		}
		port = "8000"
	} else {
		// Si no existe el archivo, variables de entorno (deploy)
		token := os.Getenv("DISCORD_TOKEN")
		channelID := os.Getenv("MONITOR_CHANNEL_ID")

		// Añade un log para verificar que las variables se están leyendo
		log.Printf("Token leído: %s", token)
		log.Printf("Channel ID leído: %s", channelID)

		cfg = &config.Config{
			Token:     token,
			ChannelID: channelID,
		}
		port = os.Getenv("PORT")
	}

	discordBot, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("Error iniciando bot: %v", err)
	}
	// Iniciar servidor HTTP
	httpServer := server.New(port)
	go httpServer.Start()

	// Iniciar bot
	if err := discordBot.Start(ctx); err != nil {
		log.Fatalf("Error iniciando bot: %v", err)
	}

	// Esperar señal de término
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Limpieza
	discordBot.Shutdown()
	httpServer.Shutdown(ctx)
}
