package service

import (
	"context"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/entity"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/pkg/logger"
)

type targetAccountService struct {
	repo   repository_contract.TargetAccountRepository
	logger logger.Logger
}

func NewTargetAccountService(repo repository_contract.TargetAccountRepository, log logger.Logger) *targetAccountService {
	return &targetAccountService{
		repo:   repo,
		logger: log,
	}
}

func (s *targetAccountService) Create(ctx context.Context, tokenContext service_contract.TokenContext, req service_contract.CreateTargetAccountRequest) (*service_contract.TargetAccountResponse, error) {
	target := &entity.TargetAccount{
		OwnerUserID: req.OwnerUserID,
		Username:    req.Username,
		Enable:      false,
		Completed:   false,
	}

	if err := s.repo.Create(ctx, target); err != nil {
		s.logger.Error(logger.Service, logger.TargetAccountService, "failed to create target account", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	s.logger.Info(logger.Service, logger.TargetAccountService, "target account created successfully", map[logger.ExtraKey]interface{}{
		logger.TargetAccountID: target.ID,
	})

	return service_contract.ToTargetAccountResponse(target), nil
}

func (s *targetAccountService) GetByID(ctx context.Context, tokenContext service_contract.TokenContext, id int64) (*service_contract.TargetAccountResponse, error) {
	target, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(logger.Service, logger.TargetAccountService, "failed to get target account by id", map[logger.ExtraKey]interface{}{
			logger.TargetAccountID: id,
			logger.ErrorMessage:    err.Error(),
		})
		return nil, err
	}

	return service_contract.ToTargetAccountResponse(target), nil
}

func (s *targetAccountService) GetByUsername(ctx context.Context, tokenContext service_contract.TokenContext, username string) (*service_contract.TargetAccountResponse, error) {
	target, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		s.logger.Error(logger.Service, logger.TargetAccountService, "failed to get target account by username", map[logger.ExtraKey]interface{}{
			logger.Username:     username,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	return service_contract.ToTargetAccountResponse(target), nil
}

func (s *targetAccountService) GetByOwnerID(ctx context.Context, tokenContext service_contract.TokenContext, ownerUserID int64, limit, offset int) ([]service_contract.TargetAccountResponse, error) {
	accounts, err := s.repo.GetByOwnerID(ctx, ownerUserID, limit, offset)
	if err != nil {
		s.logger.Error(logger.Service, logger.TargetAccountService, "failed to get target accounts by owner id", map[logger.ExtraKey]interface{}{
			logger.OwnerID:      ownerUserID,
			logger.Limit:        limit,
			logger.Offset:       offset,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	return service_contract.ToTargetAccountSliceResponse(accounts), nil
}

func (s *targetAccountService) Update(ctx context.Context, tokenContext service_contract.TokenContext, req service_contract.UpdateTargetAccountRequest) (*service_contract.TargetAccountResponse, error) {
	target := &entity.TargetAccount{
		BaseEntity:  entity.BaseEntity{ID: int64(req.ID)},
		OwnerUserID: req.OwnerUserID,
		Username:    req.Username,
		DayDuration: &req.DayDuration,
		Period:      &req.Period,
		UserLang:    req.UserLang,
		Description: &req.Description,
		Enable:      req.Enable,
		Completed:   req.Completed,
	}

	if err := s.repo.Update(ctx, target); err != nil {
		s.logger.Error(logger.Service, logger.TargetAccountService, "failed to update target account", map[logger.ExtraKey]interface{}{
			logger.TargetAccountID: req.ID,
			logger.ErrorMessage:    err.Error(),
		})
		return nil, err
	}

	s.logger.Info(logger.Service, logger.TargetAccountService, "target account updated successfully", map[logger.ExtraKey]interface{}{
		logger.TargetAccountID: req.ID,
	})

	return service_contract.ToTargetAccountResponse(target), nil
}

func (s *targetAccountService) Delete(ctx context.Context, tokenContext service_contract.TokenContext, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error(logger.Service, logger.TargetAccountService, "failed to delete target account", map[logger.ExtraKey]interface{}{
			logger.TargetAccountID: id,
			logger.ErrorMessage:    err.Error(),
		})
		return err
	}

	s.logger.Info(logger.Service, logger.TargetAccountService, "target account deleted successfully", map[logger.ExtraKey]interface{}{
		logger.TargetAccountID: id,
	})
	return nil
}

func (s *targetAccountService) CountByOwnerID(ctx context.Context, tokenContext service_contract.TokenContext, ownerUserID int64) (int64, error) {
	count, err := s.repo.CountByOwnerID(ctx, ownerUserID)
	if err != nil {
		s.logger.Error(logger.Service, logger.TargetAccountService, "failed to count target accounts by owner id", map[logger.ExtraKey]interface{}{
			logger.OwnerID:      ownerUserID,
			logger.ErrorMessage: err.Error(),
		})
		return 0, err
	}

	return count, nil
}

func (s *targetAccountService) DeleteByIDAndOwnerID(ctx context.Context, tokenContext service_contract.TokenContext, id int64, ownerUserID int64) error {
	if err := s.repo.DeleteByIDAndOwnerID(ctx, id, ownerUserID); err != nil {
		s.logger.Error(logger.Service, logger.TargetAccountService, "failed to delete target account by id and owner id", map[logger.ExtraKey]interface{}{
			logger.TargetAccountID: id,
			logger.OwnerID:         ownerUserID,
			logger.ErrorMessage:    err.Error(),
		})
		return err
	}

	s.logger.Info(logger.Service, logger.TargetAccountService, "target account deleted by owner successfully", map[logger.ExtraKey]interface{}{
		logger.TargetAccountID: id,
		logger.OwnerID:         ownerUserID,
	})
	return nil
}

func (s *targetAccountService) GetByIDAndOwnerID(ctx context.Context, tokenContext service_contract.TokenContext, id int64, ownerUserID int64) (*service_contract.TargetAccountResponse, error) {
	target, err := s.repo.GetByIDAndOwnerID(ctx, id, ownerUserID)
	if err != nil {
		s.logger.Error(logger.Service, logger.TargetAccountService, "failed to get target account by id and owner id", map[logger.ExtraKey]interface{}{
			logger.TargetAccountID: id,
			logger.OwnerID:         ownerUserID,
			logger.ErrorMessage:    err.Error(),
		})
		return nil, err
	}

	return service_contract.ToTargetAccountResponse(target), nil
}
