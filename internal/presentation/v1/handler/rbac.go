package handler

import (
	"context"
	"log"
	"strings"

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
	randomService  service_contract.RandomService
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

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        MsgShareAccountMenu,
		ReplyMarkup: ui.ShareAccessInlineKeyboard(),
	})
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send share account menu message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return
	}
}

func (h *RBACHandler) ShareAccountAccessHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := helper.GetChatID(update)
	// messageID := helper.GetMessageID(update)

	// authToken, ok := h.commonHandler.GetAuthTokenWithAccount(ctx)
	// if !ok {
	// 	return
	// }

	log.Println(update.CallbackQuery.Data)

	f := strings.TrimPrefix(update.CallbackQuery.Data, ShareAccountHandlerPrefix)
	switch f {
	case ShareAccountAccessInvite:
		_ = h.shareAccountInvite(ctx, b, chatID)
	case ShareAccountAccessList:
		_ = h.shareAccountList(ctx, b, chatID)
		return
	case ShareAccountAccessDashboard:
		_ = h.shareAccountGetBackToDashboard(ctx, b, chatID)
		return
	default:
		h.logger.Warn(logger.Handler, logger.Telegram, "unknown share account access callback action", map[logger.ExtraKey]interface{}{
			"callback_data": update.CallbackQuery.Data,
		})
	}
}

func (h *RBACHandler) shareAccountGetBackToDashboard(ctx context.Context, b *bot.Bot, chatID int64) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        MsgShareAccountMenu,
		ReplyMarkup: ui.ShareAccessInlineKeyboard(),
	})
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send share account dashboard message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return err
	}
	return nil
}

func (h *RBACHandler) shareAccountList(ctx context.Context, b *bot.Bot, chatID int64) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "This is the access list:",
	})
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send share account list message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return err
	}
	return nil
}

func (h *RBACHandler) shareAccountInvite(ctx context.Context, b *bot.Bot, chatID int64) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        MsgShareAccountInvite,
		ReplyMarkup: ui.SelectRoleInlineKeyboard(),
	})
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send share account invite message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return err
	}
	return nil
}

func (h *RBACHandler) InviteAccountAccessHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := helper.GetChatID(update)
	// messageID := helper.GetMessageID(update)

	// authToken, ok := h.commonHandler.GetAuthTokenWithAccount(ctx)
	// if !ok {
	// 	return
	// }

	log.Println(update.CallbackQuery.Data)

	f := strings.TrimPrefix(update.CallbackQuery.Data, InviteAccountHandlerPrefix)
	switch f {
	case InviteAccountAccessRoleAdmin, InviteAccountAccessRoleOwner, InviteAccountAccessRoleEditor, InviteAccountAccessRoleViewer:
		_ = h.inviteAccountRole(ctx, b, chatID, f)
	case InviteAccountAccessRoleCancel:
		_ = h.inviteAccountCancel(ctx, b, chatID)
	default:
		h.logger.Warn(logger.Handler, logger.Telegram, "unknown invite account access callback action", map[logger.ExtraKey]interface{}{
			"callback_data": update.CallbackQuery.Data,
		})
	}
}

func (h *RBACHandler) inviteAccountCancel(ctx context.Context, b *bot.Bot, chatID int64) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        MsgShareAccountMenu,
		ReplyMarkup: ui.ShareAccessInlineKeyboard(),
	})
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send cancel invite account message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return err
	}
	return nil
}

func (h *RBACHandler) inviteAccountRole(ctx context.Context, b *bot.Bot, chatID int64, role string) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   role,
	})
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send invite account role message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return err
	}
	return nil
}
