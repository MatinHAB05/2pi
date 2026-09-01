package repository_contract

import (
	"context"
	"time"
)

type ArticleDocument struct {
	ID          string            `json:"id"`
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

type ArticleSearchResult struct {
	Articles []*ArticleDocument
	Total    int
}

type ArticleESRepository interface {
	Index(ctx context.Context, article *ArticleDocument) (*string, error)
	Delete(ctx context.Context, id string) (*string, error)
	Search(ctx context.Context, query string, page, size int) (*ArticleSearchResult, error)
	Exists(ctx context.Context, id string) (*bool, error)
}
