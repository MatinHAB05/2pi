package service_contract

import (
	"context"
	"time"
)

type ArticleSearchService interface {
	IndexArticle(ctx context.Context, req *IndexArticleRequest) error
	DeleteArticle(ctx context.Context, id int64) error
	SearchArticles(ctx context.Context, req *SearchArticleQueryRequest) (*SearchArticleResponse, error)
}

type IndexArticleRequest struct {
	ID          int64  `json:"id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"  validate:"required"`
	ImageURL    string `json:"image_url"  validate:"required"`
	URL         string `json:"url"  validate:"required"`
}

type SearchArticleQueryRequest struct {
	Keyword string `json:"keyword" validate:"required"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}

type ArticleESResponse struct {
	ID          int64             `json:"id"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   time.Time         `json:"deleted_at"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	ImageURL    string            `json:"image_url"`
	URL         string            `json:"url"`
	Highlights  ArticleHighlights `json:"highlights"`
}

type ArticleHighlights struct {
	Title       []string `json:"tilte"`
	Description []string `json:"description"`
}

type SearchArticleResponse struct {
	Articles   []*ArticleESResponse `json:"articles"`
	Total      int                  `json:"total"`
	Page       int                  `json:"page"`
	Size       int                  `json:"size"`
	TotalPages int                  `json:"total_pages"`
}
