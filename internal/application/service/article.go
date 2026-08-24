package service

import (
	"context"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/entity"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/internal/infrastructure/database"
)

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
