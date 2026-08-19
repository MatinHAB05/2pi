package service_contract

import (
	"context"
	"time"

	"github.com/MatinHAB05/2pi/internal/domain/entity"
)

type UserInfoCache struct {
	Lang entity.Lang
}

type UserInfoCacheService interface {
	GetOrSyncUserInfo(
		ctx context.Context,
		tokenContext TokenContext,
		userID int64,
		ttl time.Duration,
	) (*UserInfoCache, error)

	SyncUserInfo(
		ctx context.Context,
		tokenContext TokenContext,
		userID int64,
		ttl time.Duration,
	) (*UserInfoCache, error)

	InvalidateCache(
		ctx context.Context,
		tokenContext TokenContext,
		userID int64,
	) error
}
