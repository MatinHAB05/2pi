package router

import (
	"github.com/MatinHAB05/2pi/internal/domain/entity"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
)

func UserRouter(b *bot.Bot, handlers *Handlers, services *Services, repos *Repos, middlewares *Middlewares, applogger logger.Logger) *bot.Bot {
	b.RegisterHandler(bot.HandlerTypeMessageText, "🌐 Change Language", bot.MatchTypeExact, handlers.User.ChangeLanguageMenu, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "user:fields:field:language:", bot.MatchTypePrefix, handlers.User.ChangeLanguageHandler, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))

	b.RegisterHandler(bot.HandlerTypeMessageText, "/falang", bot.MatchTypeExact, handlers.User.ChangeLanguage(entity.LangFa), bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/englang", bot.MatchTypeExact, handlers.User.ChangeLanguage(entity.LangEng), bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))

	return b
}
