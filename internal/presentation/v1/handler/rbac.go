package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/exception"
	"github.com/MatinHAB05/2pi/internal/domain/tokencontext"
	"github.com/MatinHAB05/2pi/internal/helper"
	"github.com/MatinHAB05/2pi/internal/presentation/common"
	"github.com/MatinHAB05/2pi/internal/presentation/v1/ui"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type RBACHandler struct {
	userService         service_contract.UserService
	accountService      service_contract.TargetAccountService
	rbacService         service_contract.RBACService
	shareaccountService service_contract.ShareAccountOTPService
	randomService       service_contract.RandomService
	commonHandler       *common.CommonHandler

	logger logger.Logger
}

func NewRBACHandler(
	userService service_contract.UserService,
	accountService service_contract.TargetAccountService,
	rbacService service_contract.RBACService,
	logger logger.Logger,
	commonHandler *common.CommonHandler,
	randomService service_contract.RandomService,
	shareaccountService service_contract.ShareAccountOTPService,
) RBACHandler {
	return RBACHandler{
		userService:         userService,
		accountService:      accountService,
		rbacService:         rbacService,
		logger:              logger,
		randomService:       randomService,
		commonHandler:       commonHandler,
		shareaccountService: shareaccountService,
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

	authToken, ok := h.commonHandler.GetAuthTokenWithAccount(ctx)
	if !ok {
		return
	}

	log.Println(update.CallbackQuery.Data)

	f := strings.TrimPrefix(update.CallbackQuery.Data, ShareAccountHandlerPrefix)
	switch f {
	case ShareAccountAccessInvite:
		_ = h.shareAccountInvite(ctx, b, chatID)
		return
	case ShareAccountAccessConfirmInvite:
		_ = h.shareAccountConfirmInvite(ctx, b, chatID, authToken.UserId)
		return
	case ShareAccountAccessList:
		_ = h.shareAccountList(ctx, b, chatID, *authToken.AccountID)
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
		Text:        MsgWelcome,
		ReplyMarkup: ui.MainMenuReplyKeyboard(),
	})
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send share account get back dashboard message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return err
	}
	return nil
}

func (h *RBACHandler) shareAccountList(ctx context.Context, b *bot.Bot, chatID int64, accountID int64) error {
	strAccountID := strconv.FormatInt(accountID, 10)
	res, err := h.rbacService.GetUsersForTargetAccount(ctx, service_contract.MapTokenContextToService(nil), strAccountID)
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to get users that have any access to a account ", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage:    err.Error(),
			logger.TargetAccountID: strAccountID,
		})
		return err
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprint("Who can Use this account :\n%v", res),
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

	authToken, ok := h.commonHandler.GetAuthTokenWithAccount(ctx)
	if !ok {
		return
	}

	log.Println(update.CallbackQuery.Data)

	f := strings.TrimPrefix(update.CallbackQuery.Data, InviteAccountHandlerPrefix)
	switch f {
	case InviteAccountAccessRoleAdmin, InviteAccountAccessRoleOwner, InviteAccountAccessRoleEditor, InviteAccountAccessRoleViewer:
		_ = h.inviteAccountRole(ctx, b, chatID, f, authToken)
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

func (h *RBACHandler) inviteAccountRole(ctx context.Context, b *bot.Bot, chatID int64, role string, authToken *tokencontext.AuthenticationContextToken) error {

	//TODO :  confineable len
	llleeeennnn := 8
	code, err := h.randomService.GenerateRandomBase58String(llleeeennnn)
	if err != nil {
		h.logger.Error(logger.Pkg, logger.RandomService, "failed to generate base58 random string", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err,
			"len":               llleeeennnn,
			"rand_method":       "base-58",
		})
		return err
	}

	//TODO :  confineable ttl
	h.shareaccountService.SetShareAccountOTP(ctx, service_contract.MapTokenContextToServiceJustAuth(authToken), code, service_contract.ShareAccountOTP{
		Role:          role,
		BaseAccountID: *authToken.AccountID,
		BaseUserID:    authToken.UserId,
	}, time.Hour)

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   role + ":" + code,
	})
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send invite account role message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return err
	}
	return nil
}

func (h *RBACHandler) shareAccountConfirmInvite(ctx context.Context, b *bot.Bot, chatID int64, userID int64) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   MsgEnterOneTimeShareAccountAccessCode,
	})
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send confirm share account invite input message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return err
	}
	UserStates[userID] = setConfirmShareAccountStateMap()
	return nil
}

func setConfirmShareAccountStateMap() map[string]any {
	return map[string]any{
		StateKeyStatus: StateConfirmShareAccountAccessCodeEnter,
	}
}

func getConfirmShareAccountStateMap(keyUserID int64) (state string, accountID int64, field string, err error) {
	if obj, ok := UserStates[keyUserID]; ok {
		if state, ok := obj[StateKeyStatus].(string); ok && strings.HasPrefix(state, StateEditFieldsEnter) {
			accountID = obj[StateKeyAccID].(int64)
			field = obj[StateKeyField].(string)

		}
		err = fmt.Errorf("user is not in EditAccountState")
		return
	}
	err = fmt.Errorf("not founded in user-state")
	return
}

func (h *RBACHandler) IsEnterConfirmShareAccessAccountCodeState(obj map[string]any, update *models.Update) bool {
	state, ok := obj[StateKeyStatus].(string)
	return ok && strings.HasPrefix(state, StateConfirmShareAccountAccessCodeEnter) && update.Message != nil
}

func (h *RBACHandler) handlerEnterConfirmShareAccessAccountCode(ctx context.Context, b *bot.Bot, chatID int64, userID int64, update *models.Update) error {
	code := update.Message.Text

	sh, err := h.shareaccountService.GetShareAccountOTP(ctx, service_contract.MapTokenContextToService(nil), code)
	if err != nil {
		if errors.Is(err, exception.ErrShareAccountAccessOTPCodeNotFound) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "❌ Invalid Share-Access-Account Code",
			})
		}
		h.logger.Error(logger.Handler, logger.Telegram, "failed to get share account OTP invite code", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return err

	}
	err = h.shareaccountService.InvalidateShareAccountOTP(ctx, service_contract.MapTokenContextToService(nil), code)
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to invalidate share account OTP invite code", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return err
	}

	// correct otp code
	//TODO : replace int instead of str(rbac service must in int not string base on business logic)
	ok, err := h.rbacService.AddUserRoleForTargetAccount(ctx, service_contract.MapTokenContextToService(nil), strconv.FormatInt(sh.BaseUserID, 10), strconv.FormatInt(sh.BaseAccountID, 10), sh.Role)
	if err != nil || !ok {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to update account", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
			"ok":                ok,
		})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   MsgErrUpdateAccountFailed,
		})
		return err
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf(MsgSuccessDone+":\n%+v", sh),
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        MsgWelcome,
		ReplyMarkup: ui.MainMenuReplyKeyboard(),
	})
	return nil

}
