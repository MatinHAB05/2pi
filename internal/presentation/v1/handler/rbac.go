package handler

import (
	"context"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/helper"
	"github.com/MatinHAB05/2pi/internal/presentation/common"
	"github.com/MatinHAB05/2pi/internal/presentation/v1/ui"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type RBACHandler struct {
	userService    service_contract.UserService
	accountService service_contract.TargetAccountService
	rbacService    service_contract.RBACService
	commonHandler  *common.CommonHandler

	logger logger.Logger
}

func NewRBACHandler(
	userService service_contract.UserService,
	accountService service_contract.TargetAccountService,
	rbacService service_contract.RBACService,
	logger logger.Logger,
	commonHandler *common.CommonHandler,

) RBACHandler {
	return RBACHandler{
		userService:    userService,
		accountService: accountService,
		rbacService:    rbacService,
		logger:         logger,
		commonHandler:  commonHandler,
	}
}

func (h *RBACHandler) ShareAccountAccess(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := helper.GetChatID(update)

	authToken, ok := h.commonHandler.GetAuthTokenWithAccount(ctx)
	if !ok {
		return
	}

	// user, err := h.userService.GetByID(ctx, service_contract.MapTokenContextToServiceJustAuth(authToken), authToken.UserId)
	// if err != nil {
	// 	h.logger.Error(logger.Handler, logger.Telegram, "failed to get user by id", map[logger.ExtraKey]interface{}{
	// 		logger.ErrorMessage: err.Error(),
	// 	})
	// 	return
	// }

	// h.rbacService.GetAccountsForUser()

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        MsgSwitchAccount,
		ReplyMarkup: ui.SwitchAccountInlineKeyboard([]ui.AccountItem{}, *authToken.AccountID),
	})
}
