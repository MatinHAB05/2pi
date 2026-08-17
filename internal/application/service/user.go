package service

import (
	"context"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/entity"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/pkg/logger"
)

type userService struct {
	repo   repository_contract.UserRepository
	logger logger.Logger
}

func NewUserService(repo repository_contract.UserRepository, log logger.Logger) *userService {
	return &userService{
		repo:   repo,
		logger: log,
	}
}

func (s *userService) Create(ctx context.Context, req service_contract.CreateUserRequest) (*service_contract.UserResponse, error) {
	user := &entity.User{
		BaseEntity: entity.BaseEntity{ID: int64(req.ID)},
	}

	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to create user", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	s.logger.Info(logger.Service, logger.UserService, "user created successfully", map[logger.ExtraKey]interface{}{
		logger.UserID: user.ID,
	})

	return service_contract.ToUserResponse(user), nil
}

func (s *userService) GetByID(ctx context.Context, id uint) (*service_contract.UserResponse, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to get user by id", map[logger.ExtraKey]interface{}{
			logger.UserID:       id,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	return service_contract.ToUserResponse(u), nil
}

func (s *userService) GetWithTargetAccounts(ctx context.Context, id uint) (*service_contract.UserResponse, error) {
	u, err := s.repo.GetWithTargetAccounts(ctx, id)
	if err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to get user with target accounts", map[logger.ExtraKey]interface{}{
			logger.UserID:       id,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	return service_contract.ToUserResponse(u), nil
}

func (s *userService) Update(ctx context.Context, req service_contract.UpdateUserRequest) (*service_contract.UserResponse, error) {
	user := &entity.User{
		BaseEntity: entity.BaseEntity{
			ID: int64(req.ID),
		},
	}

	if err := s.repo.Update(ctx, user); err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to update user", map[logger.ExtraKey]interface{}{
			logger.UserID:       req.ID,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	s.logger.Info(logger.Service, logger.UserService, "user updated successfully", map[logger.ExtraKey]interface{}{
		logger.UserID: req.ID,
	})

	return service_contract.ToUserResponse(user), nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to delete user", map[logger.ExtraKey]interface{}{
			logger.UserID:       id,
			logger.ErrorMessage: err.Error(),
		})
		return err
	}

	s.logger.Info(logger.Service, logger.UserService, "user deleted successfully", map[logger.ExtraKey]interface{}{
		logger.UserID: id,
	})
	return nil
}

func (s *userService) Exists(ctx context.Context, id uint) (bool, error) {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to check user existence", map[logger.ExtraKey]interface{}{
			logger.UserID:       id,
			logger.ErrorMessage: err.Error(),
		})
		return false, err
	}

	return exists, nil
}
