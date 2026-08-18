package helper

import "github.com/go-telegram/bot/models"

func GetChatID(u *models.Update) int64 {
	if u.Message != nil {
		return u.Message.Chat.ID
	}
	if u.CallbackQuery != nil && u.CallbackQuery.Message.Message != nil {
		return u.CallbackQuery.Message.Message.Chat.ID
	}
	return 0
}
