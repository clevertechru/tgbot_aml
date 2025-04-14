package handler

import (
	"testing"

	"github.com/clevertechru/tgbot_aml/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

type MockBot struct {
	sentMessages []tgbotapi.MessageConfig
}

func (b *MockBot) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	if msg, ok := c.(tgbotapi.MessageConfig); ok {
		b.sentMessages = append(b.sentMessages, msg)
	}
	return tgbotapi.Message{}, nil
}

type mockAMLService struct{}

func (m *mockAMLService) CheckAddress(address string) (*domain.CheckResult, error) {
	return &domain.CheckResult{
		IsSuspicious: false,
		RiskScore:    0.1,
		Details:      "test",
	}, nil
}

func (m *mockAMLService) CheckTransaction(fromAddress, toAddress string, amount float64) (*domain.CheckResult, error) {
	return &domain.CheckResult{
		IsSuspicious: false,
		RiskScore:    0.2,
		Details:      "test",
	}, nil
}

func TestTelegramHandler_HandleMessage(t *testing.T) {
	bot := &MockBot{}
	amlService := &mockAMLService{}
	handler := NewTelegramHandler(bot, amlService)

	// Test /start command
	msg := tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 123},
		Text: "/start",
		Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: 6},
		},
	}
	err := handler.HandleMessage(msg)
	assert.NoError(t, err)
	assert.Len(t, bot.sentMessages, 1)
	assert.Contains(t, bot.sentMessages[0].Text, "Welcome")

	// Test /check command
	msg = tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 123},
		Text: "/check test-address",
		Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: 6},
		},
	}
	err = handler.HandleMessage(msg)
	assert.NoError(t, err)
	assert.Len(t, bot.sentMessages, 2)
	assert.Contains(t, bot.sentMessages[1].Text, "test")

	// Test /checktx command
	msg = tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 123},
		Text: "/checktx from to 1.0",
		Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: 8},
		},
	}
	err = handler.HandleMessage(msg)
	assert.NoError(t, err)
	assert.Len(t, bot.sentMessages, 3)
	assert.Contains(t, bot.sentMessages[2].Text, "test")
}
