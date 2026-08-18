package middleware

import (
	"context"
	"strconv"
	"time"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/tokencontext"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Authentication(userAccCache service_contract.UserAccountCacheService, applogger logger.Logger) MiddlewareFunction {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			userID := GetUserIDFromUpdate(update)
			if userID == 0 {
				applogger.Warn(logger.General, logger.Startup, "failed to extract user_id from update", nil)
				return
			}

			var token tokencontext.AuthenticationContextToken = tokencontext.AuthenticationContextToken{
				UserId:         userID,
				AccountID:      nil,
				AccountOwnerID: nil,
			}

			account, err := userAccCache.GetOrSyncUserAccount(ctx, service_contract.MapTokenContextToService(&token), strconv.Itoa(int(userID)), "NO MATTER", time.Hour)
			if err != nil {
				applogger.Warn(logger.General, logger.Startup, "failed to get user-account from redis", map[logger.ExtraKey]interface{}{
					logger.ErrorMessage: err.Error(),
				})
				return
			}

			if account != nil {
				token.AccountID = &account.AccountID
				token.AccountOwnerID = &account.AccountOwnerID
			}

			ctx = tokencontext.SetTokenInContext(ctx, &token)

			next(ctx, b, update)
		}
	}
}

func GetUserIDFromUpdate(update *models.Update) int64 {
	if update == nil {
		return 0
	}

	switch {
	case update.Message != nil && update.Message.From != nil:
		return update.Message.From.ID
	case update.CallbackQuery != nil:
		return update.CallbackQuery.From.ID
	case update.InlineQuery != nil:
		return update.InlineQuery.From.ID
	case update.ChosenInlineResult != nil:
		return update.ChosenInlineResult.From.ID
	case update.EditedMessage != nil && update.EditedMessage.From != nil:
		return update.EditedMessage.From.ID
	case update.MyChatMember != nil:
		return update.MyChatMember.From.ID
	case update.ChatMember != nil:
		return update.ChatMember.From.ID
	case update.ChatJoinRequest != nil:
		return update.ChatJoinRequest.From.ID
	default:
		return 0
	}
}
