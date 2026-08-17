package middleware

import (
	"context"
	"fmt"

	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Recovery(applogger logger.Logger) MiddlewareFunction {
	var mf MiddlewareFunction = func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			defer func() {
				if rec := recover(); rec != nil {
					err, ok := rec.(error)
					if !ok {
						err = fmt.Errorf("internal panic")
					}
					applogger.Error(logger.Panic, logger.SubCategory(logger.MPanic), err.Error(), map[logger.ExtraKey]interface{}{
						"panic": rec,
					})
				}
			}()
			next(ctx, b, update)
		}
	}
	return mf
}
