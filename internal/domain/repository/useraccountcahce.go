package repository_contract

import (
	"context"
	"time"
)

type UserAccountCacheRepository interface {
	CreateOrReplace(ctx context.Context, userId string, ttl time.Duration) error
	CreateOrReplaceGet(ctx context.Context, userId string, ttl time.Duration) ([]string, error)
	Get(ctx context.Context, userId string) ([]string, error)
	GetSync(ctx context.Context, userId string, ttl time.Duration) ([]string, error)
	Delete(ctx context.Context, userId string) error
	Exists(ctx context.Context, userId string) (*bool, error)
}
