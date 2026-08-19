package router

import (
	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/entity"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/internal/presentation/common"
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
	Common  common.CommonHandler
	User    handler.UserHandler
	RBAC    handler.RBACHandler
}

func NewRouter(b *bot.Bot, handlers *Handlers, services *Services, repos *Repos, applogger logger.Logger) *bot.Bot {
	authMid := middleware.Authentication(services.UserAccountCache, applogger)
	clear := middleware.ClearUserState(&handlers.Common, applogger)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handlers.Basic.Start, bot.Middleware(authMid), bot.Middleware(clear))
	b.RegisterHandler(bot.HandlerTypeMessageText, "⚙️ Edit Account", bot.MatchTypeExact, handlers.Account.CompleteAccountSetup, bot.Middleware(authMid), bot.Middleware(clear))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "acc:edit:fields:field:handler:", bot.MatchTypePrefix, handlers.Account.EditAccountFields, bot.Middleware(authMid), bot.Middleware(clear))


	b.RegisterHandler(bot.HandlerTypeMessageText, "👥 Share & Access", bot.MatchTypeExact, handlers.RBAC.ShareAccountAccess, bot.Middleware(clear))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "acc:share:handler:", bot.MatchTypePrefix, handlers.RBAC.ShareAccountAccessHandler, bot.Middleware(authMid), bot.Middleware(clear))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "acc:share:invite:handler:role:", bot.MatchTypePrefix, handlers.RBAC.InviteAccountAccessHandler, bot.Middleware(authMid), bot.Middleware(clear))
	
	b.RegisterHandler(bot.HandlerTypeMessageText, "/share", bot.MatchTypeExact, handlers.RBAC.ShareAccountAccess, bot.Middleware(clear))

	b.RegisterHandler(bot.HandlerTypeMessageText, "📊 Active Account", bot.MatchTypeExact, handlers.Account.ShowActiveAccount, bot.Middleware(clear))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/account", bot.MatchTypeExact, handlers.Account.ShowActiveAccount, bot.Middleware(clear))

	b.RegisterHandler(bot.HandlerTypeMessageText, "🔄 Switch Account", bot.MatchTypeExact, handlers.Account.SwitchActiveAccount, bot.Middleware(clear))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/switch", bot.MatchTypeExact, handlers.Account.SwitchActiveAccount, bot.Middleware(clear))

	b.RegisterHandler(bot.HandlerTypeMessageText, "ℹ️ Help / Info", bot.MatchTypeExact, handlers.Basic.Help, bot.Middleware(authMid), bot.Middleware(clear))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, handlers.Basic.Help, bot.Middleware(clear))

	b.RegisterHandler(bot.HandlerTypeMessageText, "⚙️ Setting", bot.MatchTypeExact, handlers.Basic.Setting, bot.Middleware(authMid), bot.Middleware(clear))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/setting", bot.MatchTypeExact, handlers.Basic.Setting, bot.Middleware(authMid), bot.Middleware(clear))

	b.RegisterHandler(bot.HandlerTypeMessageText, "/falang", bot.MatchTypeExact, handlers.User.ChangeLanguage(entity.LangFa), bot.Middleware(authMid), bot.Middleware(clear))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/englang", bot.MatchTypeExact, handlers.User.ChangeLanguage(entity.LangEng), bot.Middleware(authMid), bot.Middleware(clear))

	return b
}
