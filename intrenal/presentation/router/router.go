package router

import (
	"github.com/MatinHAB05/reminder/intrenal/presentation/handler"
	"github.com/go-telegram/bot"
)

type Handlers struct {
	BasicHandler handler.BasicHandler
}

func NewRouter(b *bot.Bot, handlers *Handlers) *bot.Bot {

	return b
}
