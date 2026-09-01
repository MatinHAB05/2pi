package service_contract

import (
	"context"

	"github.com/MatinHAB05/2pi/internal/domain/entity"
)

type ArticleService interface {
	CreateArticle(ctx context.Context, req CreateArticleReq) (ArticleResponse, error)
	GetArticle(ctx context.Context, id int64) (ArticleResponse, error)
	ListArticles(ctx context.Context) ([]ArticleResponse, error)
	UpdateOrCreateCache(ctx context.Context) error
	LoadArticleCache(ctx context.Context) error
}
type CreateArticleReq struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	URL         string `json:"url" validate:"required,url"`
}

type ArticleResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	URL         string `json:"url"`
}

func MapToArticleResponse(article *entity.Article) ArticleResponse {
	return ArticleResponse{
		ID:          article.ID,
		Title:       article.Title,
		Description: article.Description,
		ImageURL:    article.ImageURL,
		URL:         article.URL,
	}
}
