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
		firstName string,
		lastName string,
		username string,
		reqAccountID string,
		ttl time.Duration,
	) (*UserAccountCache, error)

	SetGetDefaultAccountIfMiss(
		ctx context.Context,
		tokenContext TokenContext,
		userID int64,
		firstName string,
		lastName string,
		username string,
		reqAccountID string,
		ttl time.Duration,
	) (*UserAccountCache, error)

	InvalidateCache(ctx context.Context, tokenContext TokenContext, userID int64) error
	Set(ctx context.Context,
		tokenContext TokenContext,
		userID int64,
		accountID int64,
		ttl time.Duration) error
}
