package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/clevertechru/tgbot_aml/internal/config"
	"github.com/clevertechru/tgbot_aml/internal/domain"
	"github.com/clevertechru/tgbot_aml/internal/handlers"
	"github.com/clevertechru/tgbot_aml/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func main() {
	// Load config
	cfg, err := config.Load("config/config.yml")
	if err != nil {
		log.Printf("Failed to load config: %v, using defaults", err)
		cfg = config.DefaultConfig()
	}

	// Validate required environment variables
	if cfg.Telegram.Token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required but not set")
	}
	if cfg.AML.APIKey == "" {
		log.Fatal("AML_API_KEY is required but not set")
	}

	// Initialize logger
	var logger *zap.Logger
	if cfg.Logging.File != "" {
		file, err := os.OpenFile(cfg.Logging.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Printf("Failed to close log file: %v", err)
			}
		}()

		config := zap.NewProductionConfig()
		config.OutputPaths = []string{cfg.Logging.File}
		logger, err = config.Build()
		if err != nil {
			log.Fatalf("Failed to create logger: %v", err)
		}
	} else {
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			log.Fatalf("Failed to create logger: %v", err)
		}
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("Failed to sync logger: %v", err)
		}
	}()

	// Initialize bot
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		logger.Fatal("Failed to create bot", zap.Error(err))
	}

	// Initialize AML provider
	amlProvider := domain.NewChainabuseProvider()
	amlProvider.SetAPIKey(cfg.AML.APIKey)
	amlProvider.SetBaseURL(cfg.AML.BaseURL)

	// Initialize services
	amlService := services.NewAMLService(amlProvider)

	// Initialize handlers
	handler := handlers.NewHandler(bot, amlService, logger)

	// Set up update config
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	// Get updates channel
	updates := bot.GetUpdatesChan(updateConfig)

	// Set up context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start handling updates
	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			if err := handler.HandleMessage(ctx, update.Message); err != nil {
				logger.Error("Failed to handle message",
					zap.Error(err),
					zap.Int64("chat_id", update.Message.Chat.ID),
					zap.String("text", update.Message.Text),
				)
			}
		}
	}()

	logger.Info("Bot started", zap.String("username", bot.Self.UserName))

	// Wait for shutdown signal
	<-sigChan
	logger.Info("Shutting down...")
	cancel()
}
