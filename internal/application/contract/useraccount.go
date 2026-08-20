package service_contract

import (
	"context"
	"time"

	"github.com/MatinHAB05/2pi/internal/domain/entity"
)

type UserAccountCache struct {
	AccountID      int64
	AccountOwnerID int64
	Role           entity.Role
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
		lang string,
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
		lang string,
		ttl time.Duration,
	) (*UserAccountCache, error)

	InvalidateCache(ctx context.Context, tokenContext TokenContext, userID int64) error
	Set(ctx context.Context,
		tokenContext TokenContext,
		userID int64,
		accountID int64,
		role string,
		ttl time.Duration) error
}
