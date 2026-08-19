package service

import (
	"context"
	"errors"
	"log"
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
	userID int64,
	reqAccountID string,
	ttl time.Duration,
) (*service_contract.UserAccountCache, error) {
	userIDStr := strconv.FormatInt(userID, 10)

	// 1. Try to read from Redis
	cached, err := s.cacheRepo.Get(ctx, userIDStr)
	if err != nil {
		s.logger.Error(logger.Service, logger.CacheService, "failed to read user account from cache", map[logger.ExtraKey]interface{}{
			logger.UserID:       userID,
			logger.ErrorMessage: err.Error(),
		})
	} else if cached != nil {
		return &service_contract.UserAccountCache{
			AccountID:      cached.AccountID,
			AccountOwnerID: cached.AccountOwnerID,
		}, nil
	}

	// 2. Cache miss -> Sync logic
	return s.SyncUserAccount(ctx, tokenContext, userID, reqAccountID, ttl)
}

func (s *userAccountCacheService) SyncUserAccount(
	ctx context.Context,
	tokenContext service_contract.TokenContext,
	userID int64,
	reqAccountID string,
	ttl time.Duration,
) (*service_contract.UserAccountCache, error) {
	userIDStr := strconv.FormatInt(userID, 10)

	// Fetch user with accounts
	user, err := s.userRepo.GetWithTargetAccounts(ctx, userID)
	if err != nil {
		if !errors.Is(err, exception.ErrUserNotFound) {
			s.logger.Error(logger.Service, logger.CacheService, "failed to fetch user with target accounts", map[logger.ExtraKey]interface{}{
				logger.UserID:       userID,
				logger.ErrorMessage: err.Error(),
			})
			return nil, err
		}

		// Fallback: Create user and default target account if missing
		user = &entity.User{BaseEntity: entity.BaseEntity{ID: userID}, UserLang: entity.LangFa}
		if err := s.userRepo.Create(ctx, user); err != nil {
			s.logger.Error(logger.Service, logger.CacheService, "failed to create missing user", map[logger.ExtraKey]interface{}{
				logger.UserID:       userID,
				logger.ErrorMessage: err.Error(),
			})
			return nil, err
		}

		newAcc := &entity.TargetAccount{
			OwnerUserID: userID,
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

	if user != nil {
		if len(user.TargetAccounts) == 0 {
			return nil, nil
		}
		// log out scenario

		// ? FOR NOW :
		// we claim  that each user have exactly one account ownership(every user that ever start bot at least one time)
		// but each user can have multiple account with owner access(role) == so we have one real owner but multiple owner acerbity
		acc, err := s.accountRepo.GetByOwnerID(ctx, userID, 1, 0)
		if err != nil {
			return nil, err
		}
		log.Println("******** : ", acc)
		if len(acc) > 0 {
			reqAccountID = strconv.FormatInt(acc[0].ID, 10) // default user account // base on last comment!
		}
	}

	accountID, err := strconv.ParseInt(reqAccountID, 10, 64)
	if err != nil && user != nil && len(user.TargetAccounts) > 0 {
		accountID = user.TargetAccounts[0].ID
	}

	repoCacheData := &repository_contract.UserAccountCache{
		AccountID:      accountID,
		AccountOwnerID: userID,
	}

	// Invalidate & Set
	_ = s.cacheRepo.Delete(ctx, userIDStr)
	if err := s.cacheRepo.Set(ctx, userIDStr, repoCacheData, ttl); err != nil {
		s.logger.Error(logger.Service, logger.CacheService, "failed to write user account to cache", map[logger.ExtraKey]interface{}{
			logger.UserID:       userID,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	return &service_contract.UserAccountCache{
		AccountID:      accountID,
		AccountOwnerID: userID,
	}, nil
}

func (s *userAccountCacheService) InvalidateCache(ctx context.Context, tokenContext service_contract.TokenContext, userID int64) error {
	userIDStr := strconv.FormatInt(userID, 10)

	if err := s.cacheRepo.Delete(ctx, userIDStr); err != nil {
		s.logger.Error(logger.Service, logger.CacheService, "failed to invalidate user account cache", map[logger.ExtraKey]interface{}{
			logger.UserID:       userID,
			logger.ErrorMessage: err.Error(),
		})
		return err
	}
	return nil
}
