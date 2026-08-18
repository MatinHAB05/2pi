package handler

import (
	"context"
	"fmt"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/tokencontext"
	"github.com/MatinHAB05/2pi/internal/helper"
	"github.com/MatinHAB05/2pi/internal/presentation/v1/ui"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// TODO: Review and Refactor
type BasicHandler struct {
	userService service_contract.UserService
	logger      logger.Logger
}

func NewBasicHandler(
	userService service_contract.UserService,
	logger logger.Logger,
) BasicHandler {
	return BasicHandler{
		userService: userService,
		logger:      logger,
	}
}

func (bh *BasicHandler) NotFound(ctx context.Context, b *bot.Bot, update *models.Update) {
	bh.logger.Info(logger.Handler, logger.Telegram, "route not found", nil)
	// bh.Help(ctx, b, update) // TODO : Message/Chat id == panic
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: helper.GetChatID(update), Text: "unknown : /help"})
}

func (bh *BasicHandler) Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	authToken, err := tokencontext.GetTokenFromContext(ctx)
	if err != nil {
		bh.logger.Error(logger.Handler, logger.Telegram, "failed to get token from context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return
	}

	if authToken.Completed == false { // need completed
		bh.logger.Info(logger.Handler, logger.Telegram, "registered new user via start command", map[logger.ExtraKey]interface{}{
			logger.UserID: authToken.UserId,
		})

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "👋 Welcome to Period Tracker Bot!\n\nWe created your default account profile. Tracking is disabled until setup is completed.",
			// ParseMode:   models.ParseModeMarkdown,
			ReplyMarkup: ui.OnboardingInlineKeyboard(),
		})

		if err != nil {
			bh.logger.Error(logger.Handler, logger.Telegram, "failed to send welcome message", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
		}
		return
	}

	bh.logger.Info(logger.Handler, logger.Telegram, "existing user executed start command", map[logger.ExtraKey]interface{}{
		logger.UserID: authToken.UserId,
	})

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Welcome back! Select an option from the menu below.",
		ReplyMarkup: ui.MainMenuReplyKeyboard(),
	})

	if err != nil {
		bh.logger.Error(logger.Handler, logger.Telegram, "failed to send start response message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
	}
}

func (bh *BasicHandler) Help(ctx context.Context, b *bot.Bot, update *models.Update) {
	bh.logger.Info(logger.Handler, logger.Telegram, "help command executed", nil)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Need help? Here are the available commands...",
	}); err != nil {
		bh.logger.Error(logger.Handler, logger.Telegram, "failed to send help message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
	}
}

func (bh *BasicHandler) Setting(ctx context.Context, b *bot.Bot, update *models.Update) {
	authToken, err := tokencontext.GetTokenFromContext(ctx)
	if err != nil {
		bh.logger.Error(logger.Handler, logger.Telegram, "failed to get token from context in settings", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return
	}

	bh.logger.Info(logger.Handler, logger.Telegram, "settings command executed", map[logger.ExtraKey]interface{}{
		logger.UserID: authToken.UserId,
	})

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Bot settings for user ID: %d", authToken.UserId),
	}); err != nil {
		bh.logger.Error(logger.Handler, logger.Telegram, "failed to send settings message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
	}
}
