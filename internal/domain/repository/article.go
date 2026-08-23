package repository_contract

import (
	"context"

	"github.com/MatinHAB05/2pi/internal/domain/entity"
)

type ArticleRepository interface {
	Create(ctx context.Context, article *entity.Article) error
	GetByID(ctx context.Context, id int64) (*entity.Article, error)
	GetAll(ctx context.Context) ([]entity.Article, error)
}
