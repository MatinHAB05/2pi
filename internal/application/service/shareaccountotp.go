package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/pkg/logger"
)

type otpCacheService struct {
	cacheRepo repository_contract.ShareAccountOTPRepository
	logger    logger.Logger
}

func NewShareAccountOTPService(
	cacheRepo repository_contract.ShareAccountOTPRepository,
	log logger.Logger,
) service_contract.ShareAccountOTPService {
	return &otpCacheService{
		cacheRepo: cacheRepo,
		logger:    log,
	}
}

func (s *otpCacheService) SetOTP(
	ctx context.Context,
	tokenContext service_contract.TokenContext,
	otp string,
	data service_contract.ShareAccountOTP,
	ttl time.Duration,
) error {
	repoData := &repository_contract.ShareAccountOTP{
		Role:          data.Role,
		BaseAccountID: strconv.FormatInt(data.BaseAccountID, 10),
		BaseUserID:    strconv.FormatInt(data.BaseUserID, 10),
	}

	if err := s.cacheRepo.Set(ctx, otp, repoData, ttl); err != nil {
		s.logger.Error(logger.Service, logger.CacheService, "failed to set otp in cache", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return err
	}
	return nil
}

func (s *otpCacheService) GetOTP(
	ctx context.Context,
	tokenContext service_contract.TokenContext,
	otp string,
) (*service_contract.ShareAccountOTP, error) {
	cached, err := s.cacheRepo.Get(ctx, otp)
	if err != nil {
		s.logger.Error(logger.Service, logger.CacheService, "failed to read otp from cache", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}
	if cached == nil {
		return nil, nil
	}

	baseAccountID, err := strconv.ParseInt(cached.BaseAccountID, 10, 64)
	if err != nil {
		s.logger.Error(logger.Service, logger.CacheService, "failed to parse base_account_id from cache", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, fmt.Errorf("invalid base_account_id in cached otp: %w", err)
	}

	baseUserID, err := strconv.ParseInt(cached.BaseUserID, 10, 64)
	if err != nil {
		s.logger.Error(logger.Service, logger.CacheService, "failed to parse base_user_id from cache", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, fmt.Errorf("invalid base_user_id in cached otp: %w", err)
	}

	return &service_contract.ShareAccountOTP{
		Role:          cached.Role,
		BaseAccountID: baseAccountID,
		BaseUserID:    baseUserID,
	}, nil
}

func (s *otpCacheService) InvalidateOTP(
	ctx context.Context,
	tokenContext service_contract.TokenContext,
	otp string,
) error {
	if err := s.cacheRepo.Delete(ctx, otp); err != nil {
		s.logger.Error(logger.Service, logger.CacheService, "failed to invalidate otp cache", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return err
	}
	return nil
}
