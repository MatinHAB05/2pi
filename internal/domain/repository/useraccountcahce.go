package repository_contract

import (
	"context"
	"time"
)

type UserAccountCacheRepository interface {
	CreateOrReplace(ctx context.Context, userID string, ttl time.Duration) error
	CreateOrReplaceGet(ctx context.Context, userID string, ttl time.Duration) (*UserAccountCache, error)
	Get(ctx context.Context, userID string) (*UserAccountCache, error)
	Delete(ctx context.Context, userID string) error
	Exists(ctx context.Context, userID string) (bool, error)
	GetSync(ctx context.Context, userID string, ttl time.Duration) (*UserAccountCache, error)
}

type UserAccountCache struct {
	AccountID int64
	AccountOwnerID   int64
}
