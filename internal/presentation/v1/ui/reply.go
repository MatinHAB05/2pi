package ui

import "github.com/go-telegram/bot/models"

// MainMenuReplyKeyboard provides the persistent bottom menu
func MainMenuReplyKeyboard() *models.ReplyKeyboardMarkup {
	return &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: "📊 Current Account"},
				{Text: "🔄 Switch Account", Style: "danger"},
			},
			{
				{Text: "⚙️ Edit Account", Style: "success"},
				{Text: "👥 Share & Access", Style: "primary"},
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
