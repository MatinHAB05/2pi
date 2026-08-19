package router

import (
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
)

func RBACRouter(b *bot.Bot, handlers *Handlers, services *Services, repos *Repos, middlewares *Middlewares, applogger logger.Logger) *bot.Bot {

	b.RegisterHandler(bot.HandlerTypeMessageText, "👥 Share & Access", bot.MatchTypeExact, handlers.RBAC.ShareAccountAccess, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "acc:share:handler:", bot.MatchTypePrefix, handlers.RBAC.ShareAccountAccessHandler, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "acc:share:invite:handler:role:", bot.MatchTypePrefix, handlers.RBAC.InviteAccountAccessHandler, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))

	b.RegisterHandler(bot.HandlerTypeMessageText, "/share", bot.MatchTypeExact, handlers.RBAC.ShareAccountAccess, bot.Middleware(middlewares.Authentication), bot.Middleware(middlewares.Info), bot.Middleware(middlewares.ClearState))

	return b
}
