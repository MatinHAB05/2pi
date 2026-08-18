package middleware

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Recovery(applogger logger.Logger) MiddlewareFunction {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			defer func() {
				if rec := recover(); rec != nil {
					stackTrace := string(debug.Stack())

					var panicMsg string
					if err, ok := rec.(error); ok {
						panicMsg = err.Error()
					} else {
						panicMsg = fmt.Sprintf("%v", rec)
					}

					extra := map[logger.ExtraKey]interface{}{
						logger.ErrorMessage: panicMsg,
						"stack_trace":       stackTrace,
					}

					if update != nil {
						extra["update_id"] = update.ID

						if update.Message != nil {
							extra[logger.UserID] = update.Message.From.ID
							extra["chat_id"] = update.Message.Chat.ID
							extra["text"] = update.Message.Text
						} else if update.CallbackQuery != nil {
							extra[logger.UserID] = update.CallbackQuery.From.ID
							extra["callback_data"] = update.CallbackQuery.Data
						}
					}

					applogger.Error(
						logger.Panic,
						logger.SubCategory(logger.MPanic),
						fmt.Sprintf("recovered from panic: %s", panicMsg),
						extra,
					)
				}
			}()
			next(ctx, b, update)
		}
	}
}
