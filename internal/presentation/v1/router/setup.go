package router

import (
	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
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
	UserInfoCache    service_contract.UserInfoCacheService
}

type Handlers struct {
	Basic   handler.BasicHandler
	Account handler.AccountHandler
	Common  common.CommonHandler
	User    handler.UserHandler
	RBAC    handler.RBACHandler
	Admin handler.AdminHandler
}

type Middlewares struct {
	Logger         middleware.MiddlewareFunction
	Recovery       middleware.MiddlewareFunction
	Authentication middleware.MiddlewareFunction
	ClearState     middleware.MiddlewareFunction
	Info           middleware.MiddlewareFunction
}

func SetUpRouter(b *bot.Bot, handlers *Handlers, services *Services, repos *Repos, middlewares *Middlewares, applogger logger.Logger) *bot.Bot {

	BasicRouter(b, handlers, services, repos, middlewares, applogger)
	UserRouter(b, handlers, services, repos, middlewares, applogger)
	RBACRouter(b, handlers, services, repos, middlewares, applogger)
	AccountRouter(b, handlers, services, repos, middlewares, applogger)
	AdminRouter(b, handlers, services, repos, middlewares, applogger)


	return b
}
