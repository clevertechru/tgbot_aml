package main

import (
	"log"
	"os"

	"github.com/clevertechru/tgbot_aml/internal/domain"
	"github.com/clevertechru/tgbot_aml/internal/handler"
	"github.com/clevertechru/tgbot_aml/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Create AML service with mock provider
	mockProvider := domain.NewMockAMLService()
	amlService := service.NewAMLService(mockProvider)

	// Create telegram handler
	telegramHandler := handler.NewTelegramHandler(bot, amlService)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if err := telegramHandler.HandleMessage(*update.Message); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}
}
