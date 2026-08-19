package handler

import (
	"context"
	"fmt"
	"log"
	"strings"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/entity"
	"github.com/MatinHAB05/2pi/internal/helper"
	"github.com/MatinHAB05/2pi/internal/presentation/common"
	"github.com/MatinHAB05/2pi/internal/presentation/v1/ui"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type UserHandler struct {
	userService   service_contract.UserService
	commonHandler *common.CommonHandler

	logger logger.Logger
}

func NewUserHandler(
	userService service_contract.UserService,
	logger logger.Logger,
	commonHandler *common.CommonHandler,

) UserHandler {
	return UserHandler{
		userService:   userService,
		logger:        logger,
		commonHandler: commonHandler,
	}
}

func (h *UserHandler) ChangeLanguage(lang entity.Lang) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		authToken, ok := h.commonHandler.GetAuthToken(ctx)
		if !ok {
			return
		}

		h.logger.Info(logger.Handler, logger.Telegram, "change language command executed", map[logger.ExtraKey]interface{}{
			logger.UserID: authToken.UserId,
			"target_lang": string(lang),
		})

		tokenCtx := service_contract.MapTokenContextToServiceJustAuth(authToken)

		u, err := h.userService.Update(ctx, tokenCtx, service_contract.UpdateUserRequest{
			ID:   authToken.UserId,
			Lang: lang,
		})
		if err != nil {
			h.logger.Error(logger.Handler, logger.Telegram, "failed to update user language", map[logger.ExtraKey]interface{}{
				logger.UserID:       authToken.UserId,
				logger.ErrorMessage: err.Error(),
			})
			return
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: helper.GetChatID(update),
			Text:   fmt.Sprintf(MsgChangeLanguage, u.Lang, lang),
		})

		if err != nil {
			h.logger.Error(logger.Handler, logger.Telegram, "failed to language changed message", map[logger.ExtraKey]interface{}{
				logger.UserID:       authToken.UserId,
				logger.ErrorMessage: err.Error(),
			})
		}
	}
}

func (h *UserHandler) ChangeLanguageMenu(ctx context.Context, b *bot.Bot, update *models.Update) {
	authToken, ok := h.commonHandler.GetAuthToken(ctx)
	if !ok {
		return
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      helper.GetChatID(update),
		Text:        MsgChangeLanguageMenu,
		ReplyMarkup: ui.ChangeLanguageInlineKeyboard(),
	})

	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send change language menu message", map[logger.ExtraKey]interface{}{
			logger.UserID:       authToken.UserId,
			logger.ErrorMessage: err.Error(),
		})
	}
}

func (h *UserHandler) ChangeLanguageHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := helper.GetChatID(update)

	authToken, ok := h.commonHandler.GetAuthTokenWithAccount(ctx)
	if !ok {
		return
	}

	log.Println(update.CallbackQuery.Data)
	var lang entity.Lang
	f := strings.TrimPrefix(update.CallbackQuery.Data, UserChangeLanguageHandlerPrefix)
	switch f {

	case UserChangeLanguageHandlerFarsi:
		lang = entity.LangFa

	case UserChangeLanguageHandlerEnglish:
		lang = entity.LangEng

	case UserChangeLanguageHandlerCancel:
		h.changeLanguageCancel(ctx, b, chatID)
		return

	default:
		log.Println("WTF")
		return
	}

	u, err := h.userService.Update(ctx, service_contract.MapTokenContextToServiceJustAuth(authToken), service_contract.UpdateUserRequest{
		ID:   authToken.UserId,
		Lang: lang,
	})
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to update user language", map[logger.ExtraKey]interface{}{
			logger.UserID:       authToken.UserId,
			logger.ErrorMessage: err.Error(),
		})
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf(MsgChangeLanguage, u.Lang, lang),
	})

	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to language changed message", map[logger.ExtraKey]interface{}{
			logger.UserID:       authToken.UserId,
			logger.ErrorMessage: err.Error(),
		})
	}

}

func (h *UserHandler) changeLanguageCancel(ctx context.Context, b *bot.Bot, chatID int64) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        MsgWelcome,
		ReplyMarkup: ui.MainMenuReplyKeyboard(),
	})
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send cancel change language message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return err
	}
	return nil
}
