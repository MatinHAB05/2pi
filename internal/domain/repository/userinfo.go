package repository_contract

import (
	"context"
	"time"

	"github.com/MatinHAB05/2pi/internal/domain/entity"
)

type UserInfoCache struct {
	Lang entity.Lang
}

type UserInfoCacheRepository interface {
	Set(ctx context.Context, userID string, cache *UserInfoCache, ttl time.Duration) error
	Get(ctx context.Context, userID string) (*UserInfoCache, error)
	Delete(ctx context.Context, userID string) error
	Exists(ctx context.Context, userID string) (*bool, error)
}
