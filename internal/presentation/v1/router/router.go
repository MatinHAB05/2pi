package router

import (
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
}

type Handlers struct {
	BasicHandler handler.BasicHandler
}

func NewRouter(b *bot.Bot, handlers *Handlers, services *Services, repos *Repos, applogger logger.Logger) *bot.Bot {
	authMid := middleware.Authentication(repos.UserAccCache, applogger)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handlers.BasicHandler.Start, bot.Middleware(authMid))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, handlers.BasicHandler.Help)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/setting", bot.MatchTypeExact, handlers.BasicHandler.Setting, bot.Middleware(authMid))

	return b
}
