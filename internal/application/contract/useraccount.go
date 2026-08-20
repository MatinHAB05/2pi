package service_contract

import (
	"context"
	"time"
)

type UserAccountCache struct {
	AccountID      int64
	AccountOwnerID int64
}

type UserAccountCacheService interface {
	GetOrSetGetDefaultAccount(
		ctx context.Context,
		tokenContext TokenContext,
		userID int64,
		reqAccountID string,
		ttl time.Duration,
	) (*UserAccountCache, error)

	SetGetDefaultAccountIfMiss(
		ctx context.Context,
		tokenContext TokenContext,
		userID int64,
		reqAccountID string,
		ttl time.Duration,
	) (*UserAccountCache, error)

	Set(
		ctx context.Context,
		tokenContext TokenContext,
		userID int64,
		accountID int64,
		ttl time.Duration,
	) error

	InvalidateCache(
		ctx context.Context,
		tokenContext TokenContext,
		userID int64,
	) error
}
