package middleware

import (
	"context"
	"time"

	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ctxKey string

const loggerKey ctxKey = "logger"

func Logger(applogger logger.Logger) MiddlewareFunction {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			start := time.Now()

			info := extractUpdateInfo(update)

			applogger.Info(logger.General, logger.LoggerM, "telegram update received",
				map[logger.ExtraKey]interface{}{
					"update_id": update.ID,
					"type":      info.Type,
					"user_id":   info.UserID,
					"username":  info.Username,
					"chat_id":   info.ChatID,
					"payload":   info.Payload,
				},
			)

			newCtx := context.WithValue(ctx, loggerKey, applogger)

			next(newCtx, b, update)

			applogger.Info(logger.General, logger.LoggerM, "telegram update processed",
				map[logger.ExtraKey]interface{}{
					"update_id":   update.ID,
					"duration_ms": time.Since(start).Milliseconds(),
				},
			)
		}
	}
}

type updateInfo struct {
	Type     string
	UserID   int64
	Username string
	ChatID   int64
	Payload  string
}

func extractUpdateInfo(u *models.Update) updateInfo {
	info := updateInfo{Type: "unknown"}

	if u.Message != nil {
		info.Type = "message"
		info.Payload = u.Message.Text
		info.ChatID = u.Message.Chat.ID
		if u.Message.From != nil {
			info.UserID = u.Message.From.ID
			info.Username = u.Message.From.Username
		}
		return info
	}

	if u.CallbackQuery != nil {
		info.Type = "callback_query"
		info.Payload = u.CallbackQuery.Data
		if u.CallbackQuery.From.ID != 0 {
			info.UserID = u.CallbackQuery.From.ID
			info.Username = u.CallbackQuery.From.Username
		}
		if u.CallbackQuery.Message.Message != nil {
			info.ChatID = u.CallbackQuery.Message.Message.Chat.ID
		}
		return info
	}

	if u.EditedMessage != nil {
		info.Type = "edited_message"
		info.Payload = u.EditedMessage.Text
		info.ChatID = u.EditedMessage.Chat.ID
		if u.EditedMessage.From != nil {
			info.UserID = u.EditedMessage.From.ID
			info.Username = u.EditedMessage.From.Username
		}
		return info
	}

	return info
}
