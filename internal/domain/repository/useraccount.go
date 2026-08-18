package repository_contract

import (
	"context"
	"time"
)

type UserAccountCacheRepository interface {
	Set(ctx context.Context, userID string, cache *UserAccountCache, ttl time.Duration) error
	Get(ctx context.Context, userID string) (*UserAccountCache, error)
	Delete(ctx context.Context, userID string) error
	Exists(ctx context.Context, userID string) (bool, error)
}

type UserAccountCache struct {
	AccountID      int64
	AccountOwnerID int64
}
