package repository

import (
	"context"
	"errors"

	"github.com/MatinHAB05/2pi/internal/domain/entity"
	"github.com/MatinHAB05/2pi/internal/domain/exception"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"gorm.io/gorm"
)

type articleRepository struct {
	db *gorm.DB
}

// NewArticleRepository creates a new instance of repository_contract.ArticleRepository
func NewArticleRepository(db *gorm.DB) repository_contract.ArticleRepository {
	return &articleRepository{
		db: db,
	}
}

func (r *articleRepository) Create(ctx context.Context, article *entity.Article) error {
	return r.db.WithContext(ctx).Create(article).Error
}

func (r *articleRepository) GetByID(ctx context.Context, id int64) (*entity.Article, error) {
	var article entity.Article
	result := r.db.WithContext(ctx).First(&article, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, exception.ErrArticleNotFound // Return custom error
		}
		return nil, result.Error
	}
	return &article, nil
}
func (r *articleRepository) GetAll(ctx context.Context) ([]entity.Article, error) {
	var articles []entity.Article
	result := r.db.WithContext(ctx).Find(&articles)
	return articles, result.Error
}
