package helper

import (
	"github.com/MatinHAB05/2pi/internal/domain/entity"
	"github.com/go-telegram/bot/models"
)

func GetUserIDFromUpdate(update *models.Update) int64 {
	if update == nil {
		return 0
	}

	switch {
	case update.Message != nil && update.Message.From != nil:
		return update.Message.From.ID
	case update.CallbackQuery != nil:
		return update.CallbackQuery.From.ID
	case update.InlineQuery != nil:
		return update.InlineQuery.From.ID
	case update.ChosenInlineResult != nil:
		return update.ChosenInlineResult.From.ID
	case update.EditedMessage != nil && update.EditedMessage.From != nil:
		return update.EditedMessage.From.ID
	case update.MyChatMember != nil:
		return update.MyChatMember.From.ID
	case update.ChatMember != nil:
		return update.ChatMember.From.ID
	case update.ChatJoinRequest != nil:
		return update.ChatJoinRequest.From.ID
	default:
		return 0
	}
}

type UserInfo struct {
	FirstName    string
	LastName     string
	Username     string
	LanguageCode string
}

func GetUserDetailsFromUpdate(update *models.Update) UserInfo {
	if update == nil {
		return UserInfo{}
	}

	var user *models.User

	switch {
	case update.Message != nil && update.Message.From != nil:
		user = update.Message.From
	case update.CallbackQuery != nil:
		user = &update.CallbackQuery.From
	case update.EditedMessage != nil && update.EditedMessage.From != nil:
		user = update.EditedMessage.From
	case update.InlineQuery != nil:
		user = update.InlineQuery.From
	case update.ChosenInlineResult != nil:
		user = &update.ChosenInlineResult.From
	case update.MyChatMember != nil:
		user = &update.MyChatMember.From
	case update.ChatMember != nil:
		user = &update.ChatMember.From
	case update.ChatJoinRequest != nil:
		user = &update.ChatJoinRequest.From
	}

	if user == nil {
		return UserInfo{}
	}

	return UserInfo{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Username:     user.Username,
		LanguageCode: user.LanguageCode,
	}
}

func MapLangCodeToUserLang(langCode string) entity.Lang {
	switch langCode {
	case "fa":
		return entity.LangFa
	default: // "eng"
		return entity.LangEng
	}
}
