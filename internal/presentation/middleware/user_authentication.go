package middleware

import (
	"context"
	"time"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/tokencontext"
	"github.com/MatinHAB05/2pi/internal/helper"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Authentication(userAccCache service_contract.UserAccountCacheService, applogger logger.Logger) MiddlewareFunction {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			userID := helper.GetUserIDFromUpdate(update)
			userDetails := helper.GetUserDetailsFromUpdate(update)
			if userID == 0 {
				applogger.Warn(logger.General, logger.Startup, "failed to extract user_id from update", nil)
				return
			}

			var token tokencontext.AuthenticationContextToken = tokencontext.AuthenticationContextToken{
				UserId:         userID,
				AccountID:      nil,
				AccountOwnerID: nil,
			}

			// TODO : configurable ttl
			account, err := userAccCache.GetOrSetGetDefaultAccount(ctx, service_contract.MapTokenContextToServiceJustAuth(&token), userID, userDetails.FirstName, userDetails.LastName, userDetails.Username, "NO MATTER", time.Hour)
			if err != nil {
				applogger.Warn(logger.General, logger.Startup, "failed to get user-account from redis", map[logger.ExtraKey]interface{}{
					logger.ErrorMessage: err.Error(),
				})
				return
			}

			if account != nil {
				token.AccountID = &account.AccountID
				token.AccountOwnerID = &account.AccountOwnerID
				token.UserRole = string(account.Role)
			}

			ctx = tokencontext.SetAuthenticationTokenInContext(ctx, &token)

			next(ctx, b, update)
		}
	}
}
