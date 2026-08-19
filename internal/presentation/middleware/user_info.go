package middleware

import (
	"context"
	"time"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/tokencontext"
	"github.com/MatinHAB05/2pi/internal/presentation/common"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Info(userInfoCache service_contract.UserInfoCacheService, commonH *common.CommonHandler, applogger logger.Logger) MiddlewareFunction {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			authToken, ok := commonH.GetAuthToken(ctx)
			if !ok {
				return
			}

			// TODO : configurable ttl
			userinfo, err := userInfoCache.GetOrSyncUserInfo(ctx, service_contract.MapTokenContextToServiceJustAuth(authToken), authToken.UserId, time.Hour)
			if err != nil {
				applogger.Error(logger.Handler, logger.Telegram, "failed to get token from context", map[logger.ExtraKey]interface{}{
					logger.ErrorMessage: err.Error(),
				})
				return
			}

			var token tokencontext.UserInfoContextToken = tokencontext.UserInfoContextToken{
				Lang: userinfo.Lang,
			}

			ctx = tokencontext.SetInfoTokenInContext(ctx, &token)

			next(ctx, b, update)
		}
	}
}
