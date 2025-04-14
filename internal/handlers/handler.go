package handlers

import (
	"context"

	"github.com/clevertechru/tgbot_aml/internal/lang"
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

	userLang := lang.English
	if msg.From != nil && msg.From.LanguageCode == "ru" {
		userLang = lang.Russian
	}

	switch msg.Command() {
	case "start":
		return h.handleStart(msg, userLang)
	case "check":
		return h.handleCheck(ctx, msg, userLang)
	default:
		return h.handleUnknownCommand(msg, userLang)
	}
}

func (h *Handler) handleStart(msg *tgbotapi.Message, userLang lang.Language) error {
	reply := lang.Get(userLang, "welcome")
	response := tgbotapi.NewMessage(msg.Chat.ID, reply)
	_, err := h.bot.Send(response)
	return err
}

func (h *Handler) handleCheck(ctx context.Context, msg *tgbotapi.Message, userLang lang.Language) error {
	if len(msg.CommandArguments()) == 0 {
		reply := lang.Get(userLang, "check_usage")
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
		reply := lang.Get(userLang, "error_checking", err)
		response := tgbotapi.NewMessage(msg.Chat.ID, reply)
		_, err := h.bot.Send(response)
		return err
	}

	key := "result_clean"
	if result.IsSuspicious {
		key = "result_suspicious"
	}
	reply := lang.Get(userLang, key, result.RiskScore, result.Details)
	response := tgbotapi.NewMessage(msg.Chat.ID, reply)
	_, err = h.bot.Send(response)
	return err
}

func (h *Handler) handleUnknownCommand(msg *tgbotapi.Message, userLang lang.Language) error {
	reply := lang.Get(userLang, "unknown_command")
	response := tgbotapi.NewMessage(msg.Chat.ID, reply)
	_, err := h.bot.Send(response)
	return err
}
