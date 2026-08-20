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
				{Text: "🥷 Update Info Auto"},
				{Text: "🌐 Change Language"},
			},
			{
				{Text: "💛 Support Us" /*, WebApp: &models.WebAppInfo{URL: "https://github.com/MatinHAB05/2pi"}*/},
			},
		},
		ResizeKeyboard: true,
		IsPersistent:   false,
	}
}
