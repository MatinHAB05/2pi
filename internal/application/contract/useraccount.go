package service_contract

import (
	"context"
	"time"

	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
)

type UserAccountCacheService interface {
	GetOrSyncUserAccount(
		ctx context.Context,
		tokenContext TokenContext,
		userID string,
		reqAccountID string,
		ttl time.Duration,
	) (*repository_contract.UserAccountCache, error)

	SyncUserAccount(
		ctx context.Context,
		tokenContext TokenContext,
		userID string,
		reqAccountID string,
		ttl time.Duration,
	) (*repository_contract.UserAccountCache, error)

	InvalidateCache(
		ctx context.Context,
		tokenContext TokenContext,
		userID string,
	) error
}
