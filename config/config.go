package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Token     string `json:"token"`
	ChannelID string `json:"channel_id"`
}

func LoadConfig(file string) (*Config, error) {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return &config, err
}
