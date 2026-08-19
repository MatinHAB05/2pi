package router

import (
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
)

func AccountRouter(b *bot.Bot, handlers *Handlers, services *Services, repos *Repos, middlewares *Middlewares, applogger logger.Logger) *bot.Bot {
	b.RegisterHandler(bot.HandlerTypeMessageText, "⚙️ Edit Account", bot.MatchTypeExact, handlers.Account.CompleteAccountSetup, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "acc:edit:fields:field:handler:", bot.MatchTypePrefix, handlers.Account.EditAccountFields, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))

	b.RegisterHandler(bot.HandlerTypeMessageText, "📊 Active Account", bot.MatchTypeExact, handlers.Account.ShowActiveAccount, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/account", bot.MatchTypeExact, handlers.Account.ShowActiveAccount, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))

	b.RegisterHandler(bot.HandlerTypeMessageText, "🔄 Switch Account", bot.MatchTypeExact, handlers.Account.SwitchActiveAccount, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/switch", bot.MatchTypeExact, handlers.Account.SwitchActiveAccount, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))

	return b
}
