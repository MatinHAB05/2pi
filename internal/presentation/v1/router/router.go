package router

import (
	"github.com/MatinHAB05/2pi/internal/presentation/v1/handler"
	"github.com/go-telegram/bot"
)

type Handlers struct {
	BasicHandler handler.BasicHandler
}

func NewRouter(b *bot.Bot, handlers *Handlers) *bot.Bot {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handlers.BasicHandler.Start)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, handlers.BasicHandler.Help)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/setting", bot.MatchTypeExact, handlers.BasicHandler.Setting)

	return b
}
