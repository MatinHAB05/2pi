package handler

// func StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
// 	chatID := getChatID(update)
// 	if chatID == 0 {
// 		return
// 	}

// textParts := strings.Fields(update.Message.Text)
// if len(textParts) > 1 && strings.HasPrefix(textParts[1], "invite_") {
// 	token := strings.TrimPrefix(textParts[1], "invite_")
// 	invitation := getInvitationByToken(token)

// 	b.SendMessage(ctx, &bot.SendMessageParams{
// 		ChatID:      chatID,
// 		Text:        fmt.Sprintf("📩 **Account Invitation Received**\n\nYou have been invited to access account `@%s` as `%s`. Do you accept?", invitation.TargetUsername, invitation.Role),
// 		ParseMode:   models.ParseModeMarkdown,
// 		ReplyMarkup: ui.AcceptInviteInlineKeyboard(token),
// 	})
// 	return
// }

// 	if checkUserIsNew(chatID) { // if info not completed
// 		b.SendMessage(ctx, &bot.SendMessageParams{
// 			ChatID:      chatID,
// 			Text:        "👋 **Welcome to Period Tracker Bot!**\n\nWe created your default account profile. Tracking is **disabled** until setup is completed.",
// 			ParseMode:   models.ParseModeMarkdown,
// 			ReplyMarkup: ui.OnboardingInlineKeyboard(),
// 		})
// 		return
// 	}

// 	b.SendMessage(ctx, &bot.SendMessageParams{
// 		ChatID:      chatID,
// 		Text:        "Welcome back! Select an option from the menu below.",
// 		ReplyMarkup: ui.MainMenuReplyKeyboard(),
// 	})
// }
