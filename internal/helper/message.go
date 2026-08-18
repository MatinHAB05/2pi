package helper

import (
	"github.com/go-telegram/bot/models"
)

func FillRole(role string) string {
	return "<todo-for-now-empty>"
}

func FillDescription(description string) string {
	if description == "" {
		return "<empty>"
	}
	return description
}

func GetMessageID(update *models.Update) int {
	switch {
	case update.Message != nil:
		return update.Message.ID

	case update.EditedMessage != nil:
		return update.EditedMessage.ID

	case update.ChannelPost != nil:
		return update.ChannelPost.ID

	case update.EditedChannelPost != nil:
		return update.EditedChannelPost.ID

	case update.CallbackQuery != nil && update.CallbackQuery.Message.Message != nil:
		return update.CallbackQuery.Message.Message.ID

	default:
		return 0
	}
}
