package service

import (
	"context"
	"encoding/json"
	"os"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/entity"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/internal/infrastructure/database"
)

const cacheFilePath = "./articles_cache.json"

type articleService struct {
	repo       repository_contract.ArticleRepository
	trxManager database.TrxManager
}

func NewArticleService(
	repo repository_contract.ArticleRepository,
	trxManager database.TrxManager,
) service_contract.ArticleService {
	return &articleService{
		repo:       repo,
		trxManager: trxManager,
	}
}

func (s *articleService) CreateArticle(ctx context.Context, req service_contract.CreateArticleReq) (service_contract.ArticleResponse, error) {
	article := &entity.Article{
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		URL:         req.URL,
	}

	// 2. Save using Repository
	if err := s.repo.Create(ctx, article); err != nil {
		return service_contract.ArticleResponse{}, err
	}

	// 3. Map saved Model back to DTO
	return service_contract.MapToArticleResponse(article), nil
}

func (s *articleService) GetArticle(ctx context.Context, id int64) (service_contract.ArticleResponse, error) {
	article, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return service_contract.ArticleResponse{}, err
	}

	return service_contract.MapToArticleResponse(article), nil
}

func (s *articleService) ListArticles(ctx context.Context) ([]service_contract.ArticleResponse, error) {
	articles, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var res []service_contract.ArticleResponse
	for _, a := range articles {
		res = append(res, service_contract.MapToArticleResponse(&a))
	}

	return res, nil
}

func (s *articleService) UpdateOrCreateCache(ctx context.Context) error {
	// Skip if the cache file already exists
	if _, err := os.Stat(cacheFilePath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	// Retrieve all articles from repository to populate cache
	articles, err := s.repo.GetAll(ctx)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(articles, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFilePath, data, 0644)
}

func (s *articleService) LoadArticleCache(ctx context.Context) error {
	// Skip if the cache file does not exist
	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	data, err := os.ReadFile(cacheFilePath)
	if err != nil {
		return err
	}

	var articles []entity.Article
	if err := json.Unmarshal(data, &articles); err != nil {
		return err
	}

	for _, article := range articles {
		item := article
		if err := s.repo.Create(ctx, &item); err != nil {
			return err
		}
	}

	return nil
}
