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
	"github.com/MatinHAB05/2pi/internal/presentation/common"
	"github.com/MatinHAB05/2pi/internal/presentation/middleware"
	"github.com/MatinHAB05/2pi/internal/presentation/v1/handler"
	"github.com/MatinHAB05/2pi/internal/presentation/v1/router"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/MatinHAB05/2pi/pkg/tellogger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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
	cfg := config.Load()
	fmt.Println(cfg.String())

	// --- LOGGER INITIALIZATION ---
	telLogger, closeLF := tellogger.NewLogger(false)
	defer closeLF()

	appLogger := logger.NewLogger(logger.Config{
		Logger:   cfg.Environment.Logger.Logger,
		FilePath: cfg.Environment.Logger.FilePath,
		Encoding: cfg.Environment.Logger.Encoding,
		Level:    cfg.Environment.Logger.Level,
	})
	appLogger.Info(logger.General, logger.Startup, "starting application bootstrapping", nil)

	// --- DATABASES CONNECTIONS ---
	pgDB := database.NewPostgresDatabase(&cfg.Environment.DataBase, &cfg.Constant.Database)
	redisClient := database.NewRedisDatabase(&cfg.Environment.Redis, &cfg.Constant.Redis)

	// --- STATICS ---
	// statics, err := static.InitStaticFiles()
	// if err != nil {
	// 	appLogger.Fatal(logger.General, logger.Startup, "failed to initialize static files", map[logger.ExtraKey]interface{}{
	// 		logger.ErrorMessage: err.Error(),
	// 	})
	// }

	// --- EMAIL SENDER ---
	// emailSender := mail.NewMailer(mail.EmailConfig{
	// 	From:     cfg.Environment.Email.TwoPiEmail,
	// 	Password: cfg.Environment.Email.TwoPiEmailAppPassword,
	// 	SMTPHost: cfg.Environment.Email.SMTPHost,
	// 	SMTPPort: cfg.Environment.Email.SMTPPort,
	// })

	// casbin
	rbacEnf := rbac.NewCasbin(&cfg.Environment.Casbin, pgDB)

	// repos
	userRepo := repository.NewUserRepository(pgDB)
	userinfocacheRepo := repository.NewUserInfoCacheRepository(redisClient)
	rbacRepo := repository.NewRBACRepository(pgDB, rbacEnf)
	targetaccountRepo := repository.NewTargetAccountRepository(pgDB)
	useracccahceRepo := repository.NewUserAccountCacheRepository(redisClient)
	shareccountRepo := repository.NewShareAccountOTPCacheRepository(redisClient)
	rs := router.Repos{
		UserAccCache: useracccahceRepo,
	}

	// seed
	// rbacSeeder := seed.NewRBACSeeder(rbacRepo)
	// err := rbacSeeder.SeedAdminUser(ctx, "123456", "123456")
	// err = rbacSeeder.SeedPermissions(ctx)

	// services
	randomSrv := service.NewRandomService()
	userinfoSrv := service.NewUserInfoCacheService(userinfocacheRepo, userRepo, appLogger)
	userSrv := service.NewUserService(userRepo, appLogger, userinfoSrv)
	rbacSrv := service.NewRBACService(rbacRepo, appLogger)
	targetaccSrv := service.NewTargetAccountService(targetaccountRepo, appLogger)
	useraccountSrv := service.NewUserAccountCacheService(useracccahceRepo, userRepo, targetaccountRepo, appLogger, rbacSrv)
	shareaccountSrv := service.NewShareAccountOTPService(shareccountRepo, appLogger)
	ss := router.Services{
		UserAccountCache: useraccountSrv,
		UserInfoCache:    userinfoSrv,
	}

	// register handlers
	commonHandler := common.NewCommonHandler(appLogger)
	accountHandler := handler.NewAccountHandler(userSrv, targetaccSrv, rbacSrv, useraccountSrv, appLogger, &commonHandler)
	rbacHandler := handler.NewRBACHandler(userSrv, targetaccSrv, rbacSrv, appLogger, &commonHandler, randomSrv, shareaccountSrv)
	basicHandler := handler.NewBasicHandler(userSrv, targetaccSrv, &accountHandler, &rbacHandler, &commonHandler, appLogger)
	userHandler := handler.NewUserHandler(userSrv, appLogger, &commonHandler)
	hs := router.Handlers{
		Basic:   basicHandler,
		Account: accountHandler,
		Common:  commonHandler,
		User:    userHandler,
		RBAC:    rbacHandler,
	}

	// register middlewares
	loggerMid := middleware.Logger(appLogger)
	recoveryMid := middleware.Recovery(&cfg.Environment.DebugModeOptions, appLogger)
	authMid := middleware.Authentication(useraccountSrv, appLogger)
	clearstateMid := middleware.ClearUserState(&commonHandler, appLogger)
	infoMid := middleware.Info(userinfoSrv, &commonHandler, appLogger)
	mhs := router.Middlewares{
		Logger:         loggerMid,
		Recovery:       recoveryMid,
		Authentication: authMid,
		ClearState:     clearstateMid,
		Info:           infoMid,
	}

	// bot
	opts := []bot.Option{
		bot.WithDebugHandler(bot.DebugHandler(telLogger)),
		bot.WithDefaultHandler(basicHandler.NotFound),
		bot.WithMiddlewares(
			bot.Middleware(loggerMid),
			bot.Middleware(recoveryMid),
		),
	}
	if cfg.Environment.DebugModeOptions.Flag {
		opts = append(opts, bot.WithDebug())
	}

	b, err := bot.New(cfg.Environment.BotToken.Token, opts...)
	if err != nil {
		log.Fatalf("ERR: Failed to create bot instance: %v", err)
	}

	// register router
	router.SetUpRouter(b, &hs, &ss, &rs, &mhs, appLogger)

	// command list
	// Set Bot Commands for Telegram Menu Button
	b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{Command: "start", Description: "Start bot / Main menu"},
			{Command: "help", Description: "How to use bot"},
			{Command: "account", Description: "View current current target account"},
			{Command: "switch", Description: "Switch current target account"},
			{Command: "share", Description: "Manage access & permissions"},
			{Command: "settings", Description: "Edit target account details"},
			{Command: "falang", Description: "Change Language To Farsi"},
			{Command: "englang", Description: "Change Language To English"},
		},
	})

	// start
	fmt.Println("Bot is running...")
	b.Start(ctx)
	fmt.Println("Bot is stopped...")
}
