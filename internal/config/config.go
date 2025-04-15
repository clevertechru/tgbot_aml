package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

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
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
}

func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Replace environment variables
	re := regexp.MustCompile(`\${([^}]+)}`)
	configStr := re.ReplaceAllStringFunc(string(data), func(match string) string {
		envVar := strings.Trim(match, "${}")
		return os.Getenv(envVar)
	})

	var cfg Config
	if err := yaml.Unmarshal([]byte(configStr), &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	fmt.Printf("Server started on port %d\n", cfg.Server.Port)
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
		Server: struct {
			Port int `yaml:"port"`
		}{
			Port: 8080,
		},
	}
}
