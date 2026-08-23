package handler

import (
	"context"

	"github.com/MatinHAB05/2pi/internal/helper"
	"github.com/MatinHAB05/2pi/internal/infrastructure/scraper"
	"github.com/MatinHAB05/2pi/internal/presentation/common"
	"github.com/MatinHAB05/2pi/internal/presentation/v1/ui"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type AdminHandler struct {
	scrppers      scraper.Scrapers
	commonHandler common.CommonHandler
	logger        logger.Logger
}

func NewAdminHandler(
	scrppers scraper.Scrapers,
	logger logger.Logger,
	commonHandler common.CommonHandler,

) AdminHandler {
	return AdminHandler{
		logger:        logger,
		scrppers:      scrppers,
		commonHandler: commonHandler,
	}
}

func (h *AdminHandler) UpdateArticles(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := helper.GetChatID(update)

	authToken, ok := h.commonHandler.GetAuthToken(ctx)
	if !ok {
		h.logger.Error(logger.Handler, logger.Telegram, "unauthorized access to UpdateArticles", map[logger.ExtraKey]interface{}{
			"chat-id": chatID,
		})
		return
	}

	err := h.scrppers.RunAll(ctx, h.logger)

	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to run-all scrpapers", map[logger.ExtraKey]interface{}{
			"chat-id":           chatID,
			logger.ErrorMessage: err.Error(),
		})
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "failed to run-all scrpapers - see logs",
		})
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send err scappers", map[logger.ExtraKey]interface{}{
			"chat-id":           chatID,
			logger.ErrorMessage: err.Error(),
		})
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        MsgSuccessDone + ": All scrapers finished",
		ReplyMarkup: ui.MainMenuReplyKeyboard(),
	})
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send success article update message", map[logger.ExtraKey]interface{}{
			logger.UserID:       authToken.UserId,
			logger.ErrorMessage: err.Error(),
		})
	}
}
