package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/MatinHAB05/2pi/config"
	"github.com/MatinHAB05/2pi/internal/presentation/handler"
	"github.com/MatinHAB05/2pi/internal/presentation/router"
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
	// pgDB := database.NewPostgresDatabase(&config.Environment.DataBase, &config.Constant.Database)
	// redisClient := database.NewRedisDatabase(&config.Environment.Redis, &config.Constant.Redis)

	// register handlers
	basicHandler := handler.NewBasicHandler(appLogger)
	hs := router.Handlers{
		BasicHandler: basicHandler,
	}

	// bot
	opts := []bot.Option{
		bot.WithDebugHandler(bot.DebugHandler(telLogger)),
		bot.WithDefaultHandler(basicHandler.NotFound),
	}
	if config.Environment.DebugModeOptions.Flag {
		opts = append(opts, bot.WithDebug())
	}

	b, err := bot.New(config.Environment.BotToken.Token, opts...)
	if err != nil {
		log.Fatalf("ERR: Failed to create bot instance: %v", err)
	}

	// register router
	router.NewRouter(b, &hs)

	// start
	fmt.Println("Bot is running...")
	b.Start(ctx)
	fmt.Println("Bot is stopped...")
}
