package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Telegram struct {
		Token string `yaml:"token"`
	} `yaml:"telegram"`
	AML struct {
		APIKey  string `yaml:"api_key"`
		BaseURL string `yaml:"base_url"`
	} `yaml:"aml"`
	Logging struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	} `yaml:"logging"`
}

func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

func DefaultConfig() *Config {
	return &Config{
		Telegram: struct {
			Token string `yaml:"token"`
		}{
			Token: os.Getenv("TELEGRAM_BOT_TOKEN"),
		},
		AML: struct {
			APIKey  string `yaml:"api_key"`
			BaseURL string `yaml:"base_url"`
		}{
			APIKey:  os.Getenv("AML_API_KEY"),
			BaseURL: "https://api.aml-provider.com",
		},
		Logging: struct {
			Level string `yaml:"level"`
			File  string `yaml:"file"`
		}{
			Level: "info",
			File:  "bot.log",
		},
	}
}
