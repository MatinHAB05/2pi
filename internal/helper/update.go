package helper

import "github.com/go-telegram/bot/models"

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
