package handler

import (
	"context"
	"fmt"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/entity"
	"github.com/MatinHAB05/2pi/internal/helper"
	"github.com/MatinHAB05/2pi/internal/presentation/common"
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

		_, err := h.userService.Update(ctx, tokenCtx, service_contract.UpdateUserRequest{
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
			Text:   fmt.Sprintf(MsgSettingsFormat, authToken.UserId),
		})

		if err != nil {
			h.logger.Error(logger.Handler, logger.Telegram, "failed to send settings message", map[logger.ExtraKey]interface{}{
				logger.UserID:       authToken.UserId,
				logger.ErrorMessage: err.Error(),
			})
		}
	}
}
