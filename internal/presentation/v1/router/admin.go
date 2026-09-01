package router

import (
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
)

func AdminRouter(b *bot.Bot, handlers *Handlers, services *Services, repos *Repos, middlewares *Middlewares, applogger logger.Logger) *bot.Bot {
	//HACK : for test!
	b.RegisterHandler(bot.HandlerTypeMessageText, "/webscrap", bot.MatchTypeExact, handlers.Admin.UpdateArticles, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))

	return b
}
