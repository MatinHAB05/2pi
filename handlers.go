package main

// func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
// 	kb := &models.InlineKeyboardMarkup{
// 		InlineKeyboard: [][]models.InlineKeyboardButton{
// 			{
// 				{Text: "Option 1 🟢", CallbackData: "select_opt1"},
// 				{Text: "Option 2 🔵", CallbackData: "select_opt2"},
// 			},
// 		},
// 	}

// 	b.SendMessage(ctx, &bot.SendMessageParams{
// 		ChatID:      update.Message.Chat.ID,
// 		Text:        "Hello! Welcome to the test bot. Please select an option:",
// 		ReplyMarkup: kb,
// 	})
// }

// func helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
// 	b.SendMessage(ctx, &bot.SendMessageParams{
// 		ChatID: update.Message.Chat.ID,
// 		Text:   "Bot commands:\n/start - Start and view options\n/help - Show help message",
// 	})

// }

// func callbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
// 	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
// 		CallbackQueryID: update.CallbackQuery.ID,
// 		Text:            "Selection recorded!",
// 	})

// 	data := update.CallbackQuery.Data
// 	// m, _ :=
// 	b.SendMessage(ctx, &bot.SendMessageParams{
// 		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
// 		Text:   fmt.Sprintf("You clicked the button with data: `%s` 🔥", data),
// 	})

// 	// sm, _ := json.Marshal(m)
// 	// fmt.Println(string(sm))
// }


// 	b.SendMessage(ctx, &bot.SendMessageParams{
// 		ChatID: update.Message.Chat.ID,
// 		Text:   "Command or message not recognized. Please use /help for available options.",
// 	})
// }
