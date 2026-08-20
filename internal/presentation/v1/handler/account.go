package handler

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/tokencontext"
	"github.com/MatinHAB05/2pi/internal/helper"
	"github.com/MatinHAB05/2pi/internal/presentation/common"
	"github.com/MatinHAB05/2pi/internal/presentation/v1/ui"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type AccountHandler struct {
	userService    service_contract.UserService
	accountService service_contract.TargetAccountService
	rbacService    service_contract.RBACService
	commonHandler  *common.CommonHandler

	useraccountService service_contract.UserAccountCacheService

	logger logger.Logger
}

func NewAccountHandler(
	userService service_contract.UserService,
	accountService service_contract.TargetAccountService,
	rbacService service_contract.RBACService,
	useraccountService service_contract.UserAccountCacheService,
	logger logger.Logger,
	commonHandler *common.CommonHandler,

) AccountHandler {
	return AccountHandler{
		userService:        userService,
		accountService:     accountService,
		rbacService:        rbacService,
		useraccountService: useraccountService,
		logger:             logger,
		commonHandler:      commonHandler,
	}
}

func (h *AccountHandler) CompleteAccountSetup(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := helper.GetChatID(update)

	authToken, ok := h.commonHandler.GetAuthTokenWithAccount(ctx)
	if !ok {
		return
	}
	acc, err := h.accountService.GetByID(ctx, service_contract.MapTokenContextToServiceJustAuth(authToken), *authToken.AccountID)
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to get account by id", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        MsgEditAccountSettings,
		ReplyMarkup: ui.EditAccountInlineKeyboard(acc.Enable),
	})
}

func (h *AccountHandler) EditAccountFields(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := helper.GetChatID(update)
	messageID := helper.GetMessageID(update)

	authToken, ok := h.commonHandler.GetAuthTokenWithAccount(ctx)
	if !ok {
		return
	}

	log.Println(update.CallbackQuery.Data)
	text := ""
	f := strings.TrimPrefix(update.CallbackQuery.Data, AccountFieldHandlerPrefix)
	switch f {

	case AccountFieldDayDuration:
		text = MsgEnterNewDayDuration

	case AccountFieldPeriod:
		text = MsgEnterNewPeriod

	case AccountFieldDescription:
		text = MsgEnterNewDescription

	case AccountFieldToggleStatus:
		h.editFieldToggleEnableStatus(ctx, b, chatID, authToken)
		return
	case AccountFieldDashboard:
		h.editFieldGetBackToDashboard(ctx, b, chatID)
		return
	default:
		log.Println("WTF")
	}

	UserStates[authToken.UserId] = setEditAccountStateMap(StateEditFieldsEnter, *authToken.AccountID, f)

	if text != "" {
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    chatID,
			MessageID: messageID,
			Text:      text,
		})
	}
}

func setEditAccountStateMap(state string, AccountID int64, field string) map[string]any {
	return map[string]any{
		StateKeyStatus: state,
		StateKeyAccID:  AccountID,
		StateKeyField:  field,
	}
}

func getEditAccountStateMap(keyUserID int64) (state string, accountID int64, field string, err error) {
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

func (h *AccountHandler) editAccountFieldDayDuration(value string, ctx context.Context, b *bot.Bot, chatID int64, account *service_contract.UpdateTargetAccountRequest) error {
	dayDuration, err := strconv.Atoi(value)
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   MsgErrInvalidDayDuration,
		})
		return err
	}
	account.DayDuration = dayDuration
	return nil
}

func (h *AccountHandler) editFieldPeroid(value string, ctx context.Context, b *bot.Bot, chatID int64, account *service_contract.UpdateTargetAccountRequest) error {
	per, err := strconv.Atoi(value)
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   MsgErrInvalidPeriod,
		})
		return err
	}
	account.Period = per
	return nil
}

func (h *AccountHandler) editFieldToggleEnableStatus(ctx context.Context, b *bot.Bot, chatID int64, authToken *tokencontext.AuthenticationContextToken) error {
	acc, err := h.accountService.GetByID(ctx, service_contract.MapTokenContextToServiceJustAuth(authToken), *authToken.AccountID)
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to get account", map[logger.ExtraKey]interface{}{})
		return err
	}
	log.Println(acc.Enable)
	h.accountService.UpdateStatus(ctx, service_contract.MapTokenContextToServiceJustAuth(authToken), *authToken.AccountID, !acc.Enable)
	log.Println(!acc.Enable)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   MsgSuccessDone,
	})
	return nil
}

func (h *AccountHandler) editFieldGetBackToDashboard(ctx context.Context, b *bot.Bot, chatID int64) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        MsgWelcome,
		ReplyMarkup: ui.MainMenuReplyKeyboard(),
	})

	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send welcome message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return err
	}
	return nil
}

func (h *AccountHandler) IsEditAccountFieldState(obj map[string]any, update *models.Update) bool {
	state, ok := obj[StateKeyStatus].(string)
	return ok && strings.HasPrefix(state, StateEditFieldsEnter) && update.Message != nil
}

func (h *AccountHandler) handleEditAccountFieldState(ctx context.Context, b *bot.Bot, chatID int64, userID int64, update *models.Update) error {
	_, accID, field, err := getEditAccountStateMap(userID)
	value := update.Message.Text

	account := service_contract.UpdateTargetAccountRequest{ID: accID}
	switch field {
	case AccountFieldDayDuration:
		err := h.editAccountFieldDayDuration(value, ctx, b, chatID, &account)
		if err != nil {
			return err
		}
	case AccountFieldPeriod:
		err := h.editFieldPeroid(value, ctx, b, chatID, &account)
		if err != nil {
			return err
		}

	case AccountFieldDescription:
		account.Description = value
	}

	acc, err := h.accountService.Update(ctx, service_contract.MapTokenContextToServiceJustAuth(&tokencontext.AuthenticationContextToken{}), account)
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to update account", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   MsgErrUpdateAccountFailed,
		})
		return err
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   MsgSuccessDone,
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        MsgEditAccountSettings,
		ReplyMarkup: ui.EditAccountInlineKeyboard(acc.Enable),
	})
	return nil

}

func (h *AccountHandler) ShowCurrentAccount(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := helper.GetChatID(update)

	authToken, ok := h.commonHandler.GetAuthTokenWithAccount(ctx)
	if !ok {
		return
	}

	user, err := h.userService.GetByID(ctx, service_contract.MapTokenContextToServiceJustAuth(authToken), authToken.UserId)
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to get user by id", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return
	}

	acc, err := h.accountService.GetByID(ctx, service_contract.MapTokenContextToServiceJustAuth(authToken), *authToken.AccountID)
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to get account by id", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   ShowCurrentAccountInfo(user, acc),
		// ParseMode: models.ParseModeMarkdown,
	})
}

func (h *AccountHandler) SwitchCurrentAccount(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := helper.GetChatID(update)

	authToken, ok := h.commonHandler.GetAuthTokenWithAccount(ctx)
	if !ok {
		return
	}

	// #########
	accounts, err := h.userService.GetUserAccountsRolesByID(ctx, service_contract.MapTokenContextToService(nil), authToken.UserId)
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to get accounts-roles for user", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
			logger.UserID:       authToken.UserId,
		})
		return
	}

	h.logger.Info("", "", "", map[logger.ExtraKey]interface{}{
		"accounts": accounts,
	})

	items := ui.MapTargetAccountsToAccountItems(accounts)

	h.logger.Info("", "", "", map[logger.ExtraKey]interface{}{
		"ui-items": items,
	})

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        MsgSwitchAccount + ":",
		ReplyMarkup: ui.SwitchAccountInlineKeyboard(items, *authToken.AccountID, authToken.UserRole),
	})

	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send switch current account message", map[logger.ExtraKey]interface{}{
			logger.UserID:       authToken.UserId,
			logger.ErrorMessage: err.Error(),
		})
		return
	}
}

func (h *AccountHandler) SwitchCurrentAccountHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := helper.GetChatID(update)
	// messageID := helper.GetMessageID(update)

	authToken, ok := h.commonHandler.GetAuthTokenWithAccount(ctx)
	if !ok {
		return
	}

	log.Println(update.CallbackQuery.Data)
	temp := strings.Split(strings.TrimPrefix(update.CallbackQuery.Data, SwitchCurrentAccountHandlerPrefix), ":")
	strReqAccountID, role := temp[0], temp[1]
	reqAccountID, err := strconv.ParseInt(strReqAccountID, 10, 64)
	if err != nil {
		h.logger.Error(logger.Service, logger.CacheService, "failed to parse req_account_id from query data", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage:    err.Error(),
			logger.TargetAccountID: strReqAccountID,
		})
		return
	}

	//TODO : ttl
	err = h.useraccountService.Set(ctx, service_contract.MapTokenContextToService(nil), authToken.UserId, reqAccountID, role, time.Hour)
	if err != nil {
		h.logger.Error(logger.Service, logger.CacheService, "failed to switch/set current account", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage:    err.Error(),
			logger.TargetAccountID: reqAccountID,
			logger.UserID:          authToken.UserId,
		})
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        MsgSuccessDone + "\nuser-id : " + strconv.FormatInt(authToken.UserId, 10) + "\naccount-id : " + strconv.FormatInt(reqAccountID, 10),
		ReplyMarkup: ui.MainMenuReplyKeyboard(),
	})

	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send success switch current account message", map[logger.ExtraKey]interface{}{
			logger.UserID:       authToken.UserId,
			logger.ErrorMessage: err.Error(),
		})
		return
	}
}
