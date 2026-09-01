package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/MatinHAB05/2pi/internal/domain/entity"
	"github.com/MatinHAB05/2pi/internal/domain/exception"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/internal/infrastructure/database"
	"gorm.io/gorm"
)

type articleRepository struct {
	db database.Database
}

func NewArticleRepository(db database.Database) repository_contract.ArticleRepository {
	return &articleRepository{
		db: db,
	}
}

func (r *articleRepository) Create(ctx context.Context, article *entity.Article) error {
	db := database.ExtractTrxOrDB(ctx, r.db)
	return db.GetDB().WithContext(ctx).Create(article).Error
}

func (r *articleRepository) GetByID(ctx context.Context, id int64) (*entity.Article, error) {
	var article entity.Article
	db := database.ExtractTrxOrDB(ctx, r.db)
	result := db.GetDB().WithContext(ctx).First(&article, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, exception.ErrArticleNotFound
		}
		return nil, result.Error
	}
	return &article, nil
}

func (r *articleRepository) GetAll(ctx context.Context) ([]entity.Article, error) {
	var articles []entity.Article
	db := database.ExtractTrxOrDB(ctx, r.db)
	result := db.GetDB().WithContext(ctx).Find(&articles)
	return articles, result.Error
}

func (r *articleRepository) Exists(ctx context.Context, id int64) (*bool, error) {
	var count int64
	db := database.ExtractTrxOrDB(ctx, r.db)
	err := db.GetDB().WithContext(ctx).
		Model(&entity.Article{}).
		Where("id = ?", id).
		Count(&count).Error

	if err != nil {
		return nil, fmt.Errorf("failed to check article existence by ID %d: %w", id, err)
	}
	ex := count > 0
	return &ex, nil
}
