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
	"github.com/MatinHAB05/2pi/internal/infrastructure/repository"
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
	// rbacEnf := rbac.NewCasbin(&config.Environment.Casbin, pgDB)

	// repos
	userRepo := repository.NewUserRepository(pgDB)
	// rbacRepo := repository.NewRBACRepository(pgDB, rbacEnf)
	targetaccountRepo := repository.NewTargetAccountRepository(pgDB)
	useracccahceRepo := repository.NewUserAccountCacheRepository(redisClient)
	rs := router.Repos{
		UserAccCache: useracccahceRepo,
	}

	// seed
	// rbacSeeder := seed.NewRBACSeeder(rbacRepo)
	// err := rbacSeeder.SeedAdminUser(ctx, "123456", "123456")
	// err = rbacSeeder.SeedPermissions(ctx)

	// services
	userSrv := service.NewUserService(userRepo, appLogger)
	// rbacSrv := service.NewRBACService(rbacRepo, appLogger)
	targetaccSrv := service.NewTargetAccountService(targetaccountRepo, appLogger)
	useraccountSrv := service.NewUserAccountCacheService(useracccahceRepo, userRepo, targetaccountRepo, appLogger)
	ss := router.Services{
		UserAccountCache: useraccountSrv,
	}

	// register handlers
	basicHandler := handler.NewBasicHandler(userSrv, targetaccSrv, appLogger)
	accountHandler := handler.NewAccountHandler(userSrv, targetaccSrv, appLogger)
	hs := router.Handlers{
		Basic:   basicHandler,
		Account: accountHandler,
	}

	// bot
	opts := []bot.Option{
		bot.WithDebugHandler(bot.DebugHandler(telLogger)),
		bot.WithDefaultHandler(basicHandler.NotFound),
		bot.WithMiddlewares(
			bot.Middleware(middleware.Logger(appLogger))),
		// bot.Middleware(middleware.Recovery(appLogger))),
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

	// command list
	// Set Bot Commands for Telegram Menu Button
	b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{Command: "start", Description: "Start bot / Main menu"},
			{Command: "help", Description: "How to use bot"},
			{Command: "account", Description: "View current active target account"},
			{Command: "switch", Description: "Switch active target account"},
			{Command: "share", Description: "Manage access & permissions"},
			{Command: "settings", Description: "Edit target account details"},
		},
	})

	// start
	fmt.Println("Bot is running...")
	b.Start(ctx)
	fmt.Println("Bot is stopped...")
}
