package ui

import (
	"fmt"
	"strings"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/go-telegram/bot/models"
)

// Map single Article to InlineQueryResultArticle
func MapArticleToInlineResult(id string, article *service_contract.ArticleESResponse) models.InlineQueryResult {
	// Fallback values if fields are empty
	imageURL := article.ImageURL
	if imageURL == "" {
		imageURL = "https://via.placeholder.com/150"
	}

	description := article.Description
	if description == "" {
		description = "برای مطالعه مقاله روی دکمه زیر کلیک کنید."
	}

	artil := strings.Join(article.Highlights.Title, "")
	if artil == "" {
		artil = article.Title
	}
	artdes := strings.Join(article.Highlights.Description, "")
	if artdes == "" {
		artdes = article.Description
	}

	// messageText := fmt.Sprintf("<a href=\"%s\"><b>%s</b></a>\n\n%s", article.URL, article.Title, description)

	messageText := fmt.Sprintf("<a href=\"%s\"><b>%s</b></a>\n\n%s", article.URL, artil, artdes)

	return &models.InlineQueryResultArticle{
		ID:           id,
		Title:        article.Title,
		Description:  description,
		ThumbnailURL: imageURL,
		InputMessageContent: &models.InputTextMessageContent{
			MessageText: messageText,
			ParseMode:   models.ParseModeHTML,
		},
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text: "🔗 Open Link",
						URL:  article.URL,
					},
				},
			},
		},
	}
}

// Map slice of Articles to []models.InlineQueryResult
func MapArticlesToInlineResults(articles []*service_contract.ArticleESResponse) []models.InlineQueryResult {
	results := make([]models.InlineQueryResult, 0, len(articles))

	for i := range articles {
		id := fmt.Sprintf("article_%d", i)
		results = append(results, MapArticleToInlineResult(id, articles[i]))
	}

	return results
}
