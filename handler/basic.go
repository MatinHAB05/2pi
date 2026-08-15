package handler

import (
	"context"

	"github.com/MatinHAB05/reminder/pkg/tellog"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type BasicHandler struct {
	logger tellog.Logger
}

func NewBasicHandler(
	logger tellog.Logger,
) BasicHandler {
	return BasicHandler{
		logger: logger,
	}
}

func (bh *BasicHandler) NotFound(ctx context.Context, b *bot.Bot, update *models.Update) {
	bh.logger("basic-handler:not found handler")
}
