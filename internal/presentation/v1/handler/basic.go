package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/tokencontext"
	"github.com/MatinHAB05/2pi/internal/helper"
	"github.com/MatinHAB05/2pi/internal/presentation/v1/ui"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var UserStates = make(map[int64]map[string]any) // userid ->  |

// TODO: Review and Refactor
type BasicHandler struct {
	userService    service_contract.UserService
	accountService service_contract.TargetAccountService

	logger logger.Logger
}

func NewBasicHandler(
	userService service_contract.UserService,
	accountService service_contract.TargetAccountService,
	logger logger.Logger,
) BasicHandler {
	return BasicHandler{
		userService:    userService,
		logger:         logger,
		accountService: accountService,
	}
}

func (bh *BasicHandler) NotFound(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := helper.GetUserIDFromUpdate(update)
	chatID := helper.GetChatID(update)
	// messageID := helper.GetMessageID(update)

	if obj, ok := UserStates[userID]; ok {
		if strings.HasPrefix(obj["state"].(string), "acc:edit:fields:field:enter") && update.Message != nil {
			acc_id := obj["acc_id"].(int64)
			field := obj["field"].(string)
			value := update.Message.Text

			account := service_contract.UpdateTargetAccountRequest{ID: acc_id}
			switch field {
			case "day_duration":
				dayDuration, err := strconv.Atoi(value)
				if err != nil {
					bh.logger.Error(logger.Handler, logger.Telegram, "", map[logger.ExtraKey]interface{}{
						logger.ErrorMessage: err.Error(),
					})
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: chatID,
						Text:   "❌ Invalid day duration value. Please enter a valid number.",
					})
					return
				}
				account.DayDuration = dayDuration

			case "period":
				per, err := strconv.Atoi(value)
				if err != nil {
					bh.logger.Error(logger.Handler, logger.Telegram, "", map[logger.ExtraKey]interface{}{
						logger.ErrorMessage: err.Error(),
					})
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: chatID,
						Text:   "❌ Invalid period value. Please enter a valid number.",
					})
					return
				}
				account.Period = per

			case "description":
				account.Description = value
			}

			acc, err := bh.accountService.Update(ctx, service_contract.MapTokenContextToService(&tokencontext.AuthenticationContextToken{}), account)
			if err != nil {
				bh.logger.Error(logger.Handler, logger.Telegram, "failed to update account", map[logger.ExtraKey]interface{}{
					logger.ErrorMessage: err.Error(),
				})
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: chatID,
					Text:   "❌ Failed to update account settings. Please try again later.",
				})
				return
			}
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "✅ Done",
			})
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      chatID,
				Text:        "⚙️ Edit Account Settings \n\nSelect a parameter to update:",
				ReplyMarkup: ui.EditAccountInlineKeyboard(acc.Enable),
			})
		}
		return
	}

	bh.logger.Info(logger.Handler, logger.Telegram, "route not found", nil)
	// bh.Help(ctx, b, update) // TODO : Message/Chat id == panic
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: helper.GetChatID(update), Text: "unknown : /help"})
}

func (bh *BasicHandler) Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	authToken, err := tokencontext.GetTokenFromContext(ctx)
	if err != nil {
		bh.logger.Error(logger.Handler, logger.Telegram, "failed to get token from context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return
	}

	bh.logger.Info(logger.Handler, logger.Telegram, "registered new user via start command", map[logger.ExtraKey]interface{}{
		logger.UserID: authToken.UserId,
	})

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "👋 Welcome to Period Tracker Bot!",
		// ParseMode:   models.ParseModeMarkdown,
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
		Text:   "Need help? Here are the available commands...",
	}); err != nil {
		bh.logger.Error(logger.Handler, logger.Telegram, "failed to send help message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
	}
}

func (bh *BasicHandler) Setting(ctx context.Context, b *bot.Bot, update *models.Update) {
	authToken, err := tokencontext.GetTokenFromContext(ctx)
	if err != nil {
		bh.logger.Error(logger.Handler, logger.Telegram, "failed to get token from context in settings", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return
	}

	bh.logger.Info(logger.Handler, logger.Telegram, "settings command executed", map[logger.ExtraKey]interface{}{
		logger.UserID: authToken.UserId,
	})

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Bot settings for user ID: %d", authToken.UserId),
	}); err != nil {
		bh.logger.Error(logger.Handler, logger.Telegram, "failed to send settings message", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
	}
}
