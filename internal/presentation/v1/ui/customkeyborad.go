package ui

import (
	"github.com/go-telegram/bot/models"
)

func StartReplyKeyboardMenu(lang string) *models.ReplyKeyboardMarkup {
	return &models.ReplyKeyboardMarkup{
		IsPersistent:   false,
		ResizeKeyboard: true,
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: "Support", Style: "danger", WebApp: &models.WebAppInfo{URL: "https://github.com/MatinHAB05/2pi"}},
				{Text: "More Info", Style: "success", WebApp: &models.WebAppInfo{URL: "https://github.com/MatinHAB05/2pi/blob/master/README.md"}},
			},
			{
				{Text: "Profile", Style: "primary"},
			},
		},
	}
}
