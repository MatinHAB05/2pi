package handler

import (
	"context"
	"fmt"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/helper"
	"github.com/MatinHAB05/2pi/internal/presentation/v1/ui"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var UserStates = make(map[int64]map[string]any)

type BasicHandler struct {
	userService    service_contract.UserService
	accountService service_contract.TargetAccountService

	accountHandler *AccountHandler
	commonHandler  *CommonHandler
	logger         logger.Logger
}

func NewBasicHandler(
	userService service_contract.UserService,
	accountService service_contract.TargetAccountService,
	accountHandler *AccountHandler,
	commonHandler *CommonHandler,
	logger logger.Logger,
) BasicHandler {
	return BasicHandler{
		userService:    userService,
		logger:         logger,
		accountHandler: accountHandler,
		accountService: accountService,
		commonHandler:  commonHandler,
	}
}

func (bh *BasicHandler) NotFound(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := helper.GetUserIDFromUpdate(update)
	chatID := helper.GetChatID(update)

	if obj, ok := UserStates[userID]; ok { // stateful scenario
		if bh.accountHandler.IsEditAccountFieldState(obj, update) {
			err := bh.accountHandler.handleEditAccountFieldState(ctx, b, chatID, userID, update)
			if err != nil {
				bh.logger.Error(logger.Handler, logger.Telegram, "return err from handleEditAccountFieldState", map[logger.ExtraKey]interface{}{
					logger.ErrorMessage: err.Error(),
				})
				return
			}
		}
	}

	bh.logger.Info(logger.Handler, logger.Telegram, "route not found", nil)
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: helper.GetChatID(update), Text: MsgUnknownCommand})
}

func (bh *BasicHandler) Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	authToken, ok := bh.commonHandler.getAuthToken(ctx)
	if !ok {
		return
	}

	bh.logger.Info(logger.Handler, logger.Telegram, "registered new user via start command", map[logger.ExtraKey]interface{}{
		logger.UserID: authToken.UserId,
	})

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        MsgWelcome,
		ReplyMarkup: ui.MainMenuReplyKeyboard(),
	})

	if err != nil {
		bh.logger.Error(logger.Handler, logger.Telegram, "failed to send welcome message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
	}
}

func (bh *BasicHandler) Help(ctx context.Context, b *bot.Bot, update *models.Update) {
	bh.logger.Info(logger.Handler, logger.Telegram, "help command executed", nil)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   MsgHelp,
	}); err != nil {
		bh.logger.Error(logger.Handler, logger.Telegram, "failed to send help message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
	}
}

func (bh *BasicHandler) Setting(ctx context.Context, b *bot.Bot, update *models.Update) {
	authToken, ok := bh.commonHandler.getAuthToken(ctx)
	if !ok {
		return
	}

	bh.logger.Info(logger.Handler, logger.Telegram, "settings command executed", map[logger.ExtraKey]interface{}{
		logger.UserID: authToken.UserId,
	})

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf(MsgSettingsFormat, authToken.UserId),
	}); err != nil {
		bh.logger.Error(logger.Handler, logger.Telegram, "failed to send settings message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
	}
}
