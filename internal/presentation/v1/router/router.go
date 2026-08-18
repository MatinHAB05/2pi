package router

import (
	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/internal/presentation/middleware"
	"github.com/MatinHAB05/2pi/internal/presentation/v1/handler"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/go-telegram/bot"
)

type Repos struct {
	UserAccCache repository_contract.UserAccountCacheRepository
}

type Services struct {
	UserAccountCache service_contract.UserAccountCacheService
}

type Handlers struct {
	Basic   handler.BasicHandler
	Account handler.AccountHandler
}

func NewRouter(b *bot.Bot, handlers *Handlers, services *Services, repos *Repos, applogger logger.Logger) *bot.Bot {
	authMid := middleware.Authentication(services.UserAccountCache, applogger)
	clear := middleware.ClearUserState(applogger)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handlers.Basic.Start, bot.Middleware(authMid), bot.Middleware(clear))
	b.RegisterHandler(bot.HandlerTypeMessageText, "⚙️ Edit Account", bot.MatchTypeExact, handlers.Account.CompleteAccountSetup, bot.Middleware(authMid), bot.Middleware(clear))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "acc:edit:fields:field:handler", bot.MatchTypePrefix, handlers.Account.EditAccountFields, bot.Middleware(authMid), bot.Middleware(clear))

	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, handlers.Basic.Help, bot.Middleware(clear))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/setting", bot.MatchTypeExact, handlers.Basic.Setting, bot.Middleware(authMid), bot.Middleware(clear))

	return b
}
