package handler

import (
	"context"
	"log"
	"strings"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/tokencontext"
	"github.com/MatinHAB05/2pi/internal/helper"
	"github.com/MatinHAB05/2pi/internal/presentation/v1/ui"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type AccountHandler struct {
	userService    service_contract.UserService
	accountService service_contract.TargetAccountService

	logger logger.Logger
}

func NewAccountHandler(
	userService service_contract.UserService,
	accountService service_contract.TargetAccountService,
	logger logger.Logger,
) AccountHandler {
	return AccountHandler{
		userService:    userService,
		accountService: accountService,
		logger:         logger,
	}
}

func (h *AccountHandler) CompleteAccountSetup(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := helper.GetChatID(update)
	// messageID := helper.GetMessageID(update)

	authToken, err := tokencontext.GetTokenFromContext(ctx)
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to get token from context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return
	}
	if authToken.AccountID == nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to get account id from context - nil value", map[logger.ExtraKey]interface{}{})
		return
	}

	acc, err := h.accountService.GetByID(ctx, service_contract.MapTokenContextToService(authToken), *authToken.AccountID)
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to get account by id", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "⚙️ Edit Account Settings \n\nSelect a parameter to update:",
		// ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: ui.EditAccountInlineKeyboard(acc.Enable),
	})

}

func (h *AccountHandler) EditAccountFields(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := helper.GetChatID(update)
	messageID := helper.GetMessageID(update)

	authToken, err := tokencontext.GetTokenFromContext(ctx)
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to get token from context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return
	}
	if authToken.AccountID == nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to get account id from context - nil value", map[logger.ExtraKey]interface{}{})
		return
	}

	log.Println(update.CallbackQuery.Data)
	text := ""
	f := strings.TrimPrefix(update.CallbackQuery.Data, "acc:edit:fields:field:handler:")
	switch f {

	case "day_duration":
		text = "Enter New duration in days :"

	case "period":
		text = "Enter New period in days :"

	case "description":
		text = "Enter New Description :"

	case "toggle_status":
		acc, err := h.accountService.GetByID(ctx, service_contract.MapTokenContextToService(authToken), *authToken.AccountID)
		if err != nil {
			h.logger.Error(logger.Handler, logger.Telegram, "failed to get account", map[logger.ExtraKey]interface{}{})
			return

		}
		log.Println(acc.Enable)
		h.accountService.UpdateStatus(ctx, service_contract.MapTokenContextToService(authToken), *authToken.AccountID, !acc.Enable)
		log.Println(!acc.Enable)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "✅ Done", // answer replyCall
		})

	case "dashboard":
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "👋 Welcome to Period Tracker Bot!\n\nWe created your default account profile. Tracking is disabled until setup is completed.",
			// ParseMode:   models.ParseModeMarkdown,
			ReplyMarkup: ui.MainMenuReplyKeyboard(),
		})

		if err != nil {
			h.logger.Error(logger.Handler, logger.Telegram, "failed to send welcome message", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
		}
		return

	default:
		log.Println("WTF")
	}
	UserStates[authToken.UserId] = map[string]any{
		"state":  "acc:edit:fields:field:enter",
		"acc_id": *authToken.AccountID,
		"field":  f,
	}
	if text != "" {
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    chatID,
			MessageID: messageID,
			Text:      text,
		})
	}
}
