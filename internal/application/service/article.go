package service

import (
	"context"
	"encoding/json"
	"os"
	"strconv"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/entity"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/internal/infrastructure/database"
	"github.com/MatinHAB05/2pi/pkg/logger"
)

const cacheFilePath = "./articles_cache.json"

// todo : update es repo!
type articleService struct {
	repo   repository_contract.ArticleRepository
	esrepo repository_contract.ArticleESRepository
	logger logger.Logger

	trxManager database.TrxManager
}

func NewArticleService(
	repo repository_contract.ArticleRepository,
	trxManager database.TrxManager,
	esrepo repository_contract.ArticleESRepository,
	logger logger.Logger,
) service_contract.ArticleService {
	return &articleService{
		repo:       repo,
		trxManager: trxManager,
		esrepo:     esrepo,
		logger:     logger,
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
	for i := range articles {
		articleIDStr := strconv.FormatInt(articles[i].ID, 10)

		dbExists, err := s.repo.Exists(ctx, articles[i].ID)
		if err != nil || dbExists == nil {
			s.logger.Error(logger.Service, logger.ArticleService, "failed to check article existence in DB", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			return err
		}

		if !*dbExists {
			if err := s.repo.Create(ctx, &articles[i]); err != nil {
				return err
			}
			s.logger.Debug(logger.Service, logger.ArticleService, "article created in DB", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: nil,
			})
		} else {
			s.logger.Debug(logger.Service, logger.ArticleService, "article already exists in DB, skipping create", nil)
		}

		esExists, err := s.esrepo.Exists(ctx, articleIDStr)
		if err != nil || esExists == nil {
			s.logger.Error(logger.Service, logger.ArticleService, "failed to check article existence in ES", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			return err
		}

		if !*esExists {
			doc := &repository_contract.ArticleDocument{
				ID:          articleIDStr,
				CreatedAt:   articles[i].CreatedAt,
				UpdatedAt:   articles[i].UpdatedAt,
				DeletedAt:   articles[i].DeletedAt.Time,
				Title:       articles[i].Title,
				Description: articles[i].Description,
				ImageURL:    articles[i].ImageURL,
				URL:         articles[i].URL,
			}

			if _, err := s.esrepo.Index(ctx, doc); err != nil {
				return err
			}
			s.logger.Debug(logger.Service, logger.ArticleService, "article indexed in ES", nil)
		} else {
			s.logger.Debug(logger.Service, logger.ArticleService, "article already exists in ES, skipping index", nil)
		}
	}

	return nil
}
