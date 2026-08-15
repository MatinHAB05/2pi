package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/MatinHAB05/reminder/config"
	"github.com/MatinHAB05/reminder/handler"
	"github.com/MatinHAB05/reminder/pkg/tellog"
	"github.com/MatinHAB05/reminder/router"
	"github.com/go-telegram/bot"
)

type FF func(string) string

type X struct {
	F FF
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// init logger
	logger, closeLF := tellog.NewLogger(false)
	defer closeLF()

	// init config
	cfg := config.LoadConfig()

	// register handlers
	basicHandler := handler.NewBasicHandler(logger)

	hs := router.Handlers{
		BasicHandler: basicHandler,
	}

	// bot
	opts := []bot.Option{
		bot.WithDebugHandler(bot.DebugHandler(logger)),
		bot.WithDefaultHandler(basicHandler.NotFound),
	}
	if cfg.Debug.Flag {
		opts = append(opts, bot.WithDebug())
	}

	b, err := bot.New(cfg.BotToken.Token, opts...)
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
