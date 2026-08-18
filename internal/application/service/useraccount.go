package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/entity"
	"github.com/MatinHAB05/2pi/internal/domain/exception"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/pkg/logger"
)

type userAccountCacheService struct {
	cacheRepo   repository_contract.UserAccountCacheRepository
	userRepo    repository_contract.UserRepository
	accountRepo repository_contract.TargetAccountRepository
	logger      logger.Logger
}

func NewUserAccountCacheService(
	cacheRepo repository_contract.UserAccountCacheRepository,
	userRepo repository_contract.UserRepository,
	accountRepo repository_contract.TargetAccountRepository,
	log logger.Logger,
) service_contract.UserAccountCacheService {
	return &userAccountCacheService{
		cacheRepo:   cacheRepo,
		userRepo:    userRepo,
		accountRepo: accountRepo,
		logger:      log,
	}
}

func (s *userAccountCacheService) GetOrSyncUserAccount(
	ctx context.Context,
	tokenContext service_contract.TokenContext,
	userID string,
	reqAccountID string,
	ttl time.Duration,
) (*repository_contract.UserAccountCache, error) {
	// 1. Try to read from Redis
	cached, err := s.cacheRepo.Get(ctx, userID)
	if err != nil {
		s.logger.Error(logger.Service, logger.CacheService, "failed to read user account from cache", map[logger.ExtraKey]interface{}{
			logger.UserID:       userID,
			logger.ErrorMessage: err.Error(),
		})
	} else if cached != nil {
		return cached, nil
	}

	// 2. Cache miss -> Sync logic
	return s.SyncUserAccount(ctx, tokenContext, userID, reqAccountID, ttl)
}

func (s *userAccountCacheService) SyncUserAccount(
	ctx context.Context,
	tokenContext service_contract.TokenContext,
	userID string,
	reqAccountID string,
	ttl time.Duration,
) (*repository_contract.UserAccountCache, error) {
	intUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id format: %w", err)
	}

	// Fetch user with accounts
	user, err := s.userRepo.GetWithTargetAccounts(ctx, intUserID)
	if err != nil {
		if !errors.Is(err, exception.ErrUserNotFound) {
			s.logger.Error(logger.Service, logger.CacheService, "failed to fetch user with target accounts", map[logger.ExtraKey]interface{}{
				logger.UserID:       userID,
				logger.ErrorMessage: err.Error(),
			})
			return nil, err
		}

		// Fallback: Create user and default target account if missing
		user = &entity.User{BaseEntity: entity.BaseEntity{ID: intUserID}}
		if err := s.userRepo.Create(ctx, user); err != nil {
			s.logger.Error(logger.Service, logger.CacheService, "failed to create missing user", map[logger.ExtraKey]interface{}{
				logger.UserID:       userID,
				logger.ErrorMessage: err.Error(),
			})
			return nil, err
		}

		newAcc := &entity.TargetAccount{
			OwnerUserID: intUserID,
			Enable:      false,
		}
		if err := s.accountRepo.Create(ctx, newAcc); err != nil {
			s.logger.Error(logger.Service, logger.CacheService, "failed to create target account for user", map[logger.ExtraKey]interface{}{
				logger.UserID:       userID,
				logger.ErrorMessage: err.Error(),
			})
			return nil, err
		}

		reqAccountID = strconv.FormatInt(newAcc.ID, 10)
		user.TargetAccounts = []entity.TargetAccount{*newAcc}
	}

	if len(user.TargetAccounts) == 0 {
		return nil, nil
	}

	accountID, err := strconv.ParseInt(reqAccountID, 10, 64)
	if err != nil {
		accountID = user.TargetAccounts[0].ID
	}

	cacheData := &repository_contract.UserAccountCache{
		AccountID:      accountID,
		AccountOwnerID: intUserID,
	}

	// Invalidate & Set
	_ = s.cacheRepo.Delete(ctx, userID)
	if err := s.cacheRepo.Set(ctx, userID, cacheData, ttl); err != nil {
		s.logger.Error(logger.Service, logger.CacheService, "failed to write user account to cache", map[logger.ExtraKey]interface{}{
			logger.UserID:       userID,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	return cacheData, nil
}

func (s *userAccountCacheService) InvalidateCache(ctx context.Context, tokenContext service_contract.TokenContext, userID string) error {
	if err := s.cacheRepo.Delete(ctx, userID); err != nil {
		s.logger.Error(logger.Service, logger.CacheService, "failed to invalidate user account cache", map[logger.ExtraKey]interface{}{
			logger.UserID:       userID,
			logger.ErrorMessage: err.Error(),
		})
		return err
	}
	return nil
}
