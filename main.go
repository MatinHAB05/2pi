package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/MatinHAB05/2pi/config"
	"github.com/MatinHAB05/2pi/internal/application/service"
	"github.com/MatinHAB05/2pi/internal/infrastructure/database"
	"github.com/MatinHAB05/2pi/internal/infrastructure/rbac"
	"github.com/MatinHAB05/2pi/internal/infrastructure/repository"
	"github.com/MatinHAB05/2pi/internal/presentation/middleware"
	"github.com/MatinHAB05/2pi/internal/presentation/v1/handler"
	"github.com/MatinHAB05/2pi/internal/presentation/v1/router"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/MatinHAB05/2pi/pkg/tellogger"
	"github.com/go-telegram/bot"
)

// comment types
// TODO
// BUG
// FIXME
// HACK
// XXX
// [ ]
// [x]
// !
// ?

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// --- LOAD CONFIG & CONST ---
	config := config.Load()
	fmt.Println(config.String())

	// --- LOGGER INITIALIZATION ---
	telLogger, closeLF := tellogger.NewLogger(false)
	defer closeLF()

	appLogger := logger.NewLogger(logger.Config{
		Logger:   config.Environment.Logger.Logger,
		FilePath: config.Environment.Logger.FilePath,
		Encoding: config.Environment.Logger.Encoding,
		Level:    config.Environment.Logger.Level,
	})
	appLogger.Info(logger.General, logger.Startup, "starting application bootstrapping", nil)

	// --- DATABASES CONNECTIONS ---
	pgDB := database.NewPostgresDatabase(&config.Environment.DataBase, &config.Constant.Database)
	redisClient := database.NewRedisDatabase(&config.Environment.Redis, &config.Constant.Redis)

	// --- STATICS ---
	// statics, err := static.InitStaticFiles()
	// if err != nil {
	// 	appLogger.Fatal(logger.General, logger.Startup, "failed to initialize static files", map[logger.ExtraKey]interface{}{
	// 		logger.ErrorMessage: err.Error(),
	// 	})
	// }

	// --- EMAIL SENDER ---
	// emailSender := mail.NewMailer(mail.EmailConfig{
	// 	From:     config.Environment.Email.TwoPiEmail,
	// 	Password: config.Environment.Email.TwoPiEmailAppPassword,
	// 	SMTPHost: config.Environment.Email.SMTPHost,
	// 	SMTPPort: config.Environment.Email.SMTPPort,
	// })

	// casbin
	rbacEnf := rbac.NewCasbin(&config.Environment.Casbin, pgDB)

	// repos
	userRepo := repository.NewUserRepository(pgDB)
	rbacRepo := repository.NewRBACRepository(pgDB, rbacEnf)
	targetaccountRepo := repository.NewTargetAccountRepository(pgDB)
	useracccahceRepo := repository.NewUserAccountCacheRepository(redisClient, userRepo)
	rs := router.Repos{
		UserAccCache: useracccahceRepo,
	}

	// services
	userSrv := service.NewUserService(userRepo, appLogger)
	rbacSrv := service.NewRBACService(rbacRepo, appLogger)
	targetaccSrv := service.NewTargetAccountService(targetaccountRepo, appLogger)
	ss := router.Services{}

	// register handlers
	basicHandler := handler.NewBasicHandler(appLogger)
	hs := router.Handlers{
		BasicHandler: basicHandler,
	}

	// bot
	opts := []bot.Option{
		bot.WithDebugHandler(bot.DebugHandler(telLogger)),
		bot.WithDefaultHandler(basicHandler.NotFound),
		bot.WithMiddlewares(
			bot.Middleware(middleware.Logger(appLogger)),
			bot.Middleware(middleware.Recovery(appLogger))),
	}
	if config.Environment.DebugModeOptions.Flag {
		opts = append(opts, bot.WithDebug())
	}

	b, err := bot.New(config.Environment.BotToken.Token, opts...)
	if err != nil {
		log.Fatalf("ERR: Failed to create bot instance: %v", err)
	}

	// register router
	router.NewRouter(b, &hs, &ss, &rs, appLogger)

	// start
	fmt.Println("Bot is running...")
	b.Start(ctx)
	fmt.Println("Bot is stopped...")
}
