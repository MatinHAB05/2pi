package handler

import (
	"context"
	"fmt"
	"strings"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/helper"
	"github.com/MatinHAB05/2pi/internal/infrastructure/scraper"
	"github.com/MatinHAB05/2pi/internal/presentation/common"
	"github.com/MatinHAB05/2pi/internal/presentation/v1/ui"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var UserStates = make(map[int64]map[string]any)

type BasicHandler struct {
	userService    service_contract.UserService
	accountService service_contract.TargetAccountService

	rbacHandler    *RBACHandler
	accountHandler *AccountHandler
	commonHandler  *common.CommonHandler
	logger         logger.Logger
}

func NewBasicHandler(
	userService service_contract.UserService,
	accountService service_contract.TargetAccountService,
	accountHandler *AccountHandler,
	rbacHandler *RBACHandler,
	commonHandler *common.CommonHandler,
	logger logger.Logger,
) BasicHandler {
	return BasicHandler{
		userService:    userService,
		logger:         logger,
		accountHandler: accountHandler,
		accountService: accountService,
		commonHandler:  commonHandler,
		rbacHandler:    rbacHandler,
	}
}

func (h *BasicHandler) NotFound(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := helper.GetUserIDFromUpdate(update)
	chatID := helper.GetChatID(update)

	// inline-query scenario
	if update.InlineQuery != nil {
		// query := update.InlineQuery.Query
		_, err := b.AnswerInlineQuery(ctx, &bot.AnswerInlineQueryParams{
			InlineQueryID: update.InlineQuery.ID,
			Results:       ui.MapArticlesToInlineResults(scraper.ScrapData[:10]),
			CacheTime:     0, //todo : for debug!
		})

		if err != nil {
			fmt.Println("Error answering inline query:", err)
		}
		return
	}
	// inline-query-resault scenario
	if update.ChosenInlineResult != nil {
		rawResID := update.ChosenInlineResult.ResultID
		res_id := strings.TrimPrefix(rawResID, "article_")
		h.logger.Info(logger.Handler, logger.Scrap, "user chose artirlce from queryinline", map[logger.ExtraKey]interface{}{
			logger.UserID:      userID,
			logger.ChatID:      chatID,
			"resault-id":       res_id,
			"resault-id-in-ui": rawResID,
		})
		return
	}
	// stateful scenario
	if obj, ok := UserStates[userID]; ok {
		switch {
		case h.accountHandler.IsEditAccountFieldState(obj, update):
			err := h.accountHandler.handleEditAccountFieldState(ctx, b, chatID, userID, update)
			if err != nil {
				h.logger.Error(logger.Handler, logger.Telegram, "return err from handleEditAccountFieldState", map[logger.ExtraKey]interface{}{
					logger.ErrorMessage: err.Error(),
				})
				return
			}
			return

		case h.rbacHandler.IsEnterConfirmShareAccessAccountCodeState(obj, update):
			err := h.rbacHandler.handlerEnterConfirmShareAccessAccountCode(ctx, b, chatID, userID, update)
			if err != nil {
				h.logger.Error(logger.Handler, logger.Telegram, "return err from handlerEnterConfirmShareAccessAccountCode", map[logger.ExtraKey]interface{}{
					logger.ErrorMessage: err.Error(),
				})
				return
			}
			return
		}
		return
	}

	h.logger.Info(logger.Handler, logger.Telegram, "route not found", nil)
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: helper.GetChatID(update), Text: MsgUnknownCommand})
}

func (h *BasicHandler) Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	authToken, ok := h.commonHandler.GetAuthToken(ctx)
	if !ok {
		return
	}

	h.logger.Info(logger.Handler, logger.Telegram, "registered new user via start command", map[logger.ExtraKey]interface{}{
		logger.UserID: authToken.UserId,
	})
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        MsgWelcome,
		ReplyMarkup: ui.MainMenuReplyKeyboard(),
	})

	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send welcome message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return
	}
}

func (h *BasicHandler) Help(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Info(logger.Handler, logger.Telegram, "help command executed", nil)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   MsgHelp,
	}); err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send help message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return
	}
}

func (h *BasicHandler) Setting(ctx context.Context, b *bot.Bot, update *models.Update) {
	authToken, ok := h.commonHandler.GetAuthToken(ctx)
	if !ok {
		return
	}
	infoToken, ok := h.commonHandler.GetUserInfoToken(ctx)
	if !ok {
		return
	}

	h.logger.Info(logger.Handler, logger.Telegram, "settings command executed", map[logger.ExtraKey]interface{}{
		logger.UserID: authToken.UserId,
	})

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf(MsgSettingsFormat, authToken.UserId, infoToken.Lang),
	}); err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send settings message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return
	}
}

func (h *BasicHandler) SupportUs(ctx context.Context, b *bot.Bot, update *models.Update) {

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "💛 Support us with give star to repo 💛",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "🫶 Visit Repo", CallbackData: "support:visit-repo", URL: "https://github.com/MatinHAB05/2pi"},
				},
			},
		},
	},
	)
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to send support us message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return
	}

}
