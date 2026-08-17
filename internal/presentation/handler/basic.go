package handler

import (
	"context"

	"github.com/MatinHAB05/reminder/internal/presentation/ui/translation"
	"github.com/MatinHAB05/reminder/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// todo : REVIEW AND REFACTOR
type BasicHandler struct {
	logger logger.Logger
}

func NewBasicHandler(
	logger logger.Logger,
) BasicHandler {
	return BasicHandler{
		logger: logger,
	}
}

func (bh *BasicHandler) NotFound(ctx context.Context, b *bot.Bot, update *models.Update) {
	bh.logger.Infof("basic-handler:not found handler")
	bh.Help(ctx, b, update)
}

func (bh *BasicHandler) Start(ctx context.Context, b *bot.Bot, update *models.Update) {

	//TODO: if user registered before
	// b.SendMessage(ctx, &bot.SendMessageParams{
	// 	ChatID: update.Message.Chat.ID,
	// 	Text:   translation.StartRegistered.ToPersian(),
	// })

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   translation.StartUnknown.ToPersian(),
	})
}

func (bh *BasicHandler) Help(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   translation.HelpMe.ToPersian() + "\n--------------\n" + translation.HelpMe.ToEng(),
	})
}

func (bh *BasicHandler) Setting(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "تنظیمات بات برای کاربر گرامی ", // shows the bot's settings for this user and suggests commands to edit them.
	})
}
