package service

import (
	"context"
	"strconv"
	"time"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/pkg/logger"
)

type userInfoCacheService struct {
	userinfoCacheRepo repository_contract.UserInfoCacheRepository
	userRepo          repository_contract.UserRepository
	logger            logger.Logger
}

func NewUserInfoCacheService(
	userinfoCacheRepo repository_contract.UserInfoCacheRepository,
	userRepo repository_contract.UserRepository,
	log logger.Logger,
) service_contract.UserInfoCacheService {
	return &userInfoCacheService{
		userinfoCacheRepo: userinfoCacheRepo,
		userRepo:          userRepo,
		logger:            log,
	}
}

func (s *userInfoCacheService) GetOrSyncUserInfo(
	ctx context.Context,
	tokenContext service_contract.TokenContext,
	userID int64,
	ttl time.Duration,
) (*service_contract.UserInfoCache, error) {
	userIDStr := strconv.FormatInt(userID, 10)

	cached, err := s.userinfoCacheRepo.Get(ctx, userIDStr)
	if err != nil {
		s.logger.Error(logger.Service, logger.CacheService, "failed to read user info from cache", map[logger.ExtraKey]interface{}{
			logger.UserID:       userID,
			logger.ErrorMessage: err.Error(),
		})
	} else if cached != nil {
		return &service_contract.UserInfoCache{
			Lang: cached.Lang,
		}, nil
	}

	return s.SyncUserInfo(ctx, tokenContext, userID, ttl)
}

func (s *userInfoCacheService) SyncUserInfo(
	ctx context.Context,
	tokenContext service_contract.TokenContext,
	userID int64,
	ttl time.Duration,
) (*service_contract.UserInfoCache, error) {
	userIDStr := strconv.FormatInt(userID, 10)

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error(logger.Service, logger.CacheService, "failed to fetch user info from repo", map[logger.ExtraKey]interface{}{
			logger.UserID:       userID,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	repoCacheData := &repository_contract.UserInfoCache{
		Lang: user.UserLang,
	}

	_ = s.userinfoCacheRepo.Delete(ctx, userIDStr)
	if err := s.userinfoCacheRepo.Set(ctx, userIDStr, repoCacheData, ttl); err != nil {
		s.logger.Error(logger.Service, logger.CacheService, "failed to write user info to cache", map[logger.ExtraKey]interface{}{
			logger.UserID:       userID,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	return &service_contract.UserInfoCache{
		Lang: repoCacheData.Lang,
	}, nil
}

func (s *userInfoCacheService) InvalidateCache(
	ctx context.Context,
	tokenContext service_contract.TokenContext,
	userID int64,
) error {
	userIDStr := strconv.FormatInt(userID, 10)

	if err := s.userinfoCacheRepo.Delete(ctx, userIDStr); err != nil {
		s.logger.Error(logger.Service, logger.CacheService, "failed to invalidate user info cache", map[logger.ExtraKey]interface{}{
			logger.UserID:       userID,
			logger.ErrorMessage: err.Error(),
		})
		return err
	}
	return nil
}
