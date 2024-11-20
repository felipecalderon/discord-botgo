package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken     string
	MonitorChannelID string
	Port             string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env: %w", err)
	}

	cfg := &Config{
		DiscordToken:     os.Getenv("DISCORD_TOKEN"),
		MonitorChannelID: os.Getenv("MONITOR_CHANNEL_ID"),
		Port:             os.Getenv("PORT"),
	}

	if cfg.DiscordToken == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN no configurado")
	}

	if cfg.MonitorChannelID == "" {
		return nil, fmt.Errorf("MONITOR_CHANNEL_ID no configurado")
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg, nil
}
