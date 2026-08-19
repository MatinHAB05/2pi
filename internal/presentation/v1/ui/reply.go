package ui

import "github.com/go-telegram/bot/models"

// MainMenuReplyKeyboard provides the persistent bottom menu
func MainMenuReplyKeyboard() *models.ReplyKeyboardMarkup {
	return &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: "📊 Active Account"},
				{Text: "🔄 Switch Account"},
			},
			{
				{Text: "⚙️ Edit Account"},
				{Text: "👥 Share & Access"},
			},
			{
				{Text: "ℹ️ Help / Info"},
				{Text: "⚙️ Settings"},
			},
			{
				{Text: "🌐 Change Language"},
			},
		},
		ResizeKeyboard: true,
		IsPersistent:   false,
	}
}
