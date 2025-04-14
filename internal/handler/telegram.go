package handler

import (
	"fmt"
	"strings"

	"github.com/clevertechru/tgbot_aml/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotAPI defines the interface for sending messages
type BotAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

// TelegramHandler handles Telegram bot messages
type TelegramHandler struct {
	bot        BotAPI
	amlService domain.AMLService
}

// NewTelegramHandler creates a new TelegramHandler
func NewTelegramHandler(bot BotAPI, amlService domain.AMLService) *TelegramHandler {
	return &TelegramHandler{
		bot:        bot,
		amlService: amlService,
	}
}

// HandleMessage processes a Telegram message
func (h *TelegramHandler) HandleMessage(msg tgbotapi.Message) error {
	switch {
	case msg.IsCommand():
		return h.handleCommand(msg)
	default:
		return h.sendMessage(msg.Chat.ID, "Please use a command. Type /start to see available commands.")
	}
}

func (h *TelegramHandler) handleCommand(msg tgbotapi.Message) error {
	switch msg.Command() {
	case "start":
		return h.sendMessage(msg.Chat.ID, "Welcome to AML Checker Bot!\n\nAvailable commands:\n/check <address> - Check an address\n/checktx <from> <to> <amount> - Check a transaction")
	case "check":
		address := strings.TrimSpace(msg.CommandArguments())
		if address == "" {
			return h.sendMessage(msg.Chat.ID, "Please provide an address to check")
		}
		return h.checkAddress(msg.Chat.ID, address)
	case "checktx":
		args := strings.Fields(msg.CommandArguments())
		if len(args) != 3 {
			return h.sendMessage(msg.Chat.ID, "Please provide from address, to address, and amount")
		}
		return h.checkTransaction(msg.Chat.ID, args[0], args[1], args[2])
	default:
		return h.sendMessage(msg.Chat.ID, "Unknown command. Type /start to see available commands.")
	}
}

func (h *TelegramHandler) checkAddress(chatID int64, address string) error {
	result, err := h.amlService.CheckAddress(address)
	if err != nil {
		return h.sendMessage(chatID, fmt.Sprintf("Error checking address: %v", err))
	}

	message := fmt.Sprintf("Address check results:\nSuspicious: %v\nRisk Score: %.2f\nDetails: %s",
		result.IsSuspicious, result.RiskScore, result.Details)
	return h.sendMessage(chatID, message)
}

func (h *TelegramHandler) checkTransaction(chatID int64, from, to, amount string) error {
	result, err := h.amlService.CheckTransaction(from, to, 0) // TODO: Parse amount
	if err != nil {
		return h.sendMessage(chatID, fmt.Sprintf("Error checking transaction: %v", err))
	}

	message := fmt.Sprintf("Transaction check results:\nSuspicious: %v\nRisk Score: %.2f\nDetails: %s",
		result.IsSuspicious, result.RiskScore, result.Details)
	return h.sendMessage(chatID, message)
}

func (h *TelegramHandler) sendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := h.bot.Send(msg)
	return err
}
