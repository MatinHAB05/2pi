package router

import (
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func BasicRouter(b *bot.Bot, handlers *Handlers, services *Services, repos *Repos, middlewares *Middlewares, applogger logger.Logger) *bot.Bot {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handlers.Basic.Start, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))

	b.RegisterHandler(bot.HandlerTypeMessageText, "ℹ️ Help / Info", bot.MatchTypeExact, handlers.Basic.Help, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, handlers.Basic.Help, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))

	b.RegisterHandler(bot.HandlerTypeMessageText, "⚙️ Settings", bot.MatchTypeExact, handlers.Basic.Setting, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/settings", bot.MatchTypeExact, handlers.Basic.Setting, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))

	b.RegisterHandler(bot.HandlerTypeMessageText, "💛 Support Us", bot.MatchTypeExact, handlers.Basic.SupportUs, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))

	b.RegisterHandlerMatchFunc(func(u *models.Update) bool { return u.InlineQuery != nil }, handlers.Basic.ArticleSearchInlineQuery)

	return b
}
