package main

import (
	"context"
	"discord-bot/internal/bot"
	"discord-bot/internal/config"
	"discord-bot/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	// Inicializar bot
	discordBot, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("Error iniciando bot: %v", err)
	}

	// Iniciar servidor HTTP
	httpServer := server.New(cfg.Port)
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
