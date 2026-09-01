package service

import (
	"context"
	"math"
	"strconv"
	"time"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/pkg/logger"
)

type articleSearchService struct {
	esRepo repository_contract.ArticleESRepository
	logger logger.Logger
}

func NewArticleSearchService(esRepo repository_contract.ArticleESRepository, logger logger.Logger) service_contract.ArticleSearchService {
	return &articleSearchService{
		esRepo: esRepo,
		logger: logger,
	}
}

func (s *articleSearchService) IndexArticle(ctx context.Context, req *service_contract.IndexArticleRequest) error {
	doc := &repository_contract.ArticleDocument{
		ID:          strconv.FormatInt(req.ID, 10),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		URL:         req.URL,
	}

	if _, err := s.esRepo.Index(ctx, doc); err != nil {
		s.logger.Error(logger.Service, logger.ArticleService, "failed to index article in ES", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return err
	}

	s.logger.Info(logger.Service, logger.ArticleService, "article indexed successfully in ES", nil)
	return nil
}

func (s *articleSearchService) DeleteArticle(ctx context.Context, id int64) error {
	if _, err := s.esRepo.Delete(ctx, strconv.FormatInt(id, 10)); err != nil {
		s.logger.Error(logger.Service, logger.ArticleService, "failed to delete article from ES", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return err
	}

	s.logger.Info(logger.Service, logger.ArticleService, "article deleted successfully from ES", nil)
	return nil
}

func (s *articleSearchService) SearchArticles(ctx context.Context, req *service_contract.SearchArticleQueryRequest) (*service_contract.SearchArticleResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	s.logger.Debug(logger.Service, logger.ArticleService, "searching articles in ES", map[logger.ExtraKey]interface{}{
		logger.Limit:  req.Size,
		logger.Offset: (req.Page - 1) * req.Size,
	})

	result, err := s.esRepo.Search(ctx, req.Keyword, req.Page, req.Size)
	if err != nil {
		s.logger.Error(logger.Service, logger.ArticleService, "failed to search articles in ES", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	articles := make([]*service_contract.ArticleESResponse, 0, len(result.Articles))
	for _, doc := range result.Articles {
		parsedID, err := strconv.ParseInt(doc.ID, 10, 64)
		if err != nil {
			s.logger.Error(logger.Service, logger.ArticleService, "failed to parse article ID from string to int64", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			return nil, err
		}

		articles = append(articles, &service_contract.ArticleESResponse{
			ID:          parsedID,
			CreatedAt:   doc.CreatedAt,
			UpdatedAt:   doc.UpdatedAt,
			Title:       doc.Title,
			Description: doc.Description,
			ImageURL:    doc.ImageURL,
			URL:         doc.URL,
			DeletedAt:   doc.DeletedAt,
			Highlights: service_contract.ArticleHighlights{
				Title:       doc.Highlights.Title,
				Description: doc.Highlights.Description,
			},
		})
	}

	totalPages := int(math.Ceil(float64(result.Total) / float64(req.Size)))

	return &service_contract.SearchArticleResponse{
		Articles:   articles,
		Total:      result.Total,
		Page:       req.Page,
		Size:       req.Size,
		TotalPages: totalPages,
	}, nil
}
