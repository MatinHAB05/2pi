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
	rbacService service_contract.RBACService
	logger      logger.Logger
}

func NewUserAccountCacheService(
	cacheRepo repository_contract.UserAccountCacheRepository,
	userRepo repository_contract.UserRepository,
	accountRepo repository_contract.TargetAccountRepository,
	log logger.Logger,
	rbacService service_contract.RBACService,
) service_contract.UserAccountCacheService {
	return &userAccountCacheService{
		cacheRepo:   cacheRepo,
		userRepo:    userRepo,
		accountRepo: accountRepo,
		logger:      log,
		rbacService: rbacService,
	}
}

func (s *userAccountCacheService) GetOrSetGetDefaultAccount(
	ctx context.Context,
	tokenContext service_contract.TokenContext,
	userID int64,
	firstName string,
	lastName string,
	username string,
	reqAccountID string,
	lang string,
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
		accountID, err := strconv.ParseInt(cached.AccountID, 10, 64)
		if err != nil {
			s.logger.Error(logger.Service, logger.CacheService, "failed to parse account_id from cache", map[logger.ExtraKey]interface{}{
				logger.UserID:       userID,
				logger.ErrorMessage: err.Error(),
			})
			return nil, fmt.Errorf("invalid account_id in cache: %w", err)
		}

		accountOwnerID, err := strconv.ParseInt(cached.AccountOwnerID, 10, 64)
		if err != nil {
			s.logger.Error(logger.Service, logger.CacheService, "failed to parse account_owner_id from cache", map[logger.ExtraKey]interface{}{
				logger.UserID:       userID,
				logger.ErrorMessage: err.Error(),
			})
			return nil, fmt.Errorf("invalid account_owner_id in cache: %w", err)
		}

		return &service_contract.UserAccountCache{
			AccountID:      accountID,
			AccountOwnerID: accountOwnerID,
			Role:           entity.Role(cached.Role),
		}, nil
	}

	// 2. Cache miss -> Sync logic
	return s.SetGetDefaultAccountIfMiss(ctx, tokenContext, userID, firstName, lastName, username, lang, reqAccountID, ttl)
}

func (s *userAccountCacheService) SetGetDefaultAccountIfMiss(
	ctx context.Context,
	tokenContext service_contract.TokenContext,
	userID int64,
	firstName string,
	lastName string,
	username string,
	lang string,
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

		// Fallback: Create user with profile fields and default target account if missing
		user = &entity.User{
			BaseEntity: entity.BaseEntity{ID: userID},
			FirstName:  firstName,
			LastName:   lastName,
			Username:   username,
			UserLang:   entity.Lang(lang),
		}
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

		_, err := s.rbacService.AddUserRoleForTargetAccount(ctx, service_contract.MapTokenContextToService(nil), userID, newAcc.ID, string(entity.RoleOwner))
		if err != nil {
			return nil, err
		}
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
		if len(acc) > 0 {
			reqAccountID = strconv.FormatInt(acc[0].ID, 10) // default user account  // base on last comment!
		}
	}

	accountID, err := strconv.ParseInt(reqAccountID, 10, 64)
	if err != nil && user != nil && len(user.TargetAccounts) > 0 {
		accountID = user.TargetAccounts[0].ID
	}

	repoCacheData := &repository_contract.UserAccountCache{
		AccountID:      strconv.FormatInt(accountID, 10),
		AccountOwnerID: strconv.FormatInt(userID, 10),
		Role:           string(entity.RoleOwner),
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
		Role:           entity.RoleOwner,
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

func (s *userAccountCacheService) Set(ctx context.Context, tokenContext service_contract.TokenContext, userID int64, accountID int64, role string, ttl time.Duration) error {
	userIDStr := strconv.FormatInt(userID, 10)

	acc, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return err
	}

	err = s.cacheRepo.Set(ctx, userIDStr, &repository_contract.UserAccountCache{
		AccountID:      strconv.FormatInt(accountID, 10),
		AccountOwnerID: strconv.FormatInt(acc.OwnerUserID, 10),
		Role:           role,
	}, ttl)

	if err != nil {
		s.logger.Error(logger.Service, logger.CacheService, "failed to invalidate user account cache", map[logger.ExtraKey]interface{}{
			logger.UserID:       userID,
			logger.ErrorMessage: err.Error(),
		})
		return err
	}
	return nil
}
