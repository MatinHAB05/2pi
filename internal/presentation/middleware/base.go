package middleware

import (
	"github.com/go-telegram/bot"
)

type MiddlewareFunction func(next bot.HandlerFunc) bot.HandlerFunc
