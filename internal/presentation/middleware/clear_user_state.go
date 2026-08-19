package middleware

import (
	"context"

	"github.com/MatinHAB05/2pi/internal/presentation/common"
	"github.com/MatinHAB05/2pi/internal/presentation/v1/handler"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func ClearUserState(commonH *common.CommonHandler, applogger logger.Logger) MiddlewareFunction {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			authToken, ok := commonH.GetAuthToken(ctx)
			if !ok {
				return
			}

			delete(handler.UserStates, authToken.UserId)

			next(ctx, b, update)
		}
	}
}
