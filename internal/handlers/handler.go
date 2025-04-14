package handlers

import (
	"context"
	"fmt"

	"github.com/clevertechru/tgbot_aml/internal/domain"
	"github.com/clevertechru/tgbot_aml/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Handler struct {
	bot        *tgbotapi.BotAPI
	amlService *services.AMLService
	logger     *zap.Logger
}

func NewHandler(bot *tgbotapi.BotAPI, amlService *services.AMLService, logger *zap.Logger) *Handler {
	return &Handler{
		bot:        bot,
		amlService: amlService,
		logger:     logger,
	}
}

func (h *Handler) HandleMessage(ctx context.Context, msg *tgbotapi.Message) error {
	if msg == nil {
		return nil
	}

	switch msg.Command() {
	case "start":
		return h.handleStart(msg)
	case "check":
		return h.handleCheck(ctx, msg)
	default:
		return h.handleUnknownCommand(msg)
	}
}

func (h *Handler) handleStart(msg *tgbotapi.Message) error {
	reply := "Welcome to AML Bot! Use /check <address> to check an address or transaction."
	response := tgbotapi.NewMessage(msg.Chat.ID, reply)
	_, err := h.bot.Send(response)
	return err
}

func (h *Handler) handleCheck(ctx context.Context, msg *tgbotapi.Message) error {
	if len(msg.CommandArguments()) == 0 {
		reply := "Please provide an address or transaction hash to check. Usage: /check <address>"
		response := tgbotapi.NewMessage(msg.Chat.ID, reply)
		_, err := h.bot.Send(response)
		return err
	}

	target := msg.CommandArguments()
	result, err := h.amlService.CheckAddress(ctx, target)
	if err != nil {
		h.logger.Error("Failed to check address",
			zap.Error(err),
			zap.String("address", target),
		)
		reply := fmt.Sprintf("Error checking address: %v", err)
		response := tgbotapi.NewMessage(msg.Chat.ID, reply)
		_, err := h.bot.Send(response)
		return err
	}

	reply := domain.FormatAMLResult(result)
	response := tgbotapi.NewMessage(msg.Chat.ID, reply)
	_, err = h.bot.Send(response)
	return err
}

func (h *Handler) handleUnknownCommand(msg *tgbotapi.Message) error {
	reply := "Unknown command. Use /start to see available commands."
	response := tgbotapi.NewMessage(msg.Chat.ID, reply)
	_, err := h.bot.Send(response)
	return err
}
