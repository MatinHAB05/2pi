package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/internal/infrastructure/database"
	"github.com/redis/go-redis/v9"
)

type userAccountCacheRepository struct {
	client database.Cache
}

func NewUserAccountCacheRepository(client database.Cache) repository_contract.UserAccountCacheRepository {
	return &userAccountCacheRepository{
		client: client,
	}
}

func (ucr *userAccountCacheRepository) buildUserAccountCacheKey(userID string) string {
	return fmt.Sprintf("user:%s:account", userID)
}

func (ucr *userAccountCacheRepository) Set(ctx context.Context, userID string, cache *repository_contract.UserAccountCache, ttl time.Duration) error {
	key := ucr.buildUserAccountCacheKey(userID)

	fields := map[string]interface{}{
		"account_id":       cache.AccountID,
		"account_owner_id": cache.AccountOwnerID,
	}

	pipe := ucr.client.GetRDB().Pipeline()
	pipe.HSet(ctx, key, fields)
	if ttl > 0 {
		pipe.Expire(ctx, key, ttl)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("redis set user account cache error: %w", err)
	}

	return nil
}

func (ucr *userAccountCacheRepository) Get(ctx context.Context, userID string) (*repository_contract.UserAccountCache, error) {
	key := ucr.buildUserAccountCacheKey(userID)

	res, err := ucr.client.GetRDB().HGetAll(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("redis get user account error: %w", err)
	}

	if len(res) == 0 {
		return nil, nil
	}

	accountID, err := strconv.ParseInt(res["account_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid account_id in redis: %w", err)
	}

	accountOwnerID, err := strconv.ParseInt(res["account_owner_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid account_owner_id in redis: %w", err)
	}

	return &repository_contract.UserAccountCache{
		AccountID:      accountID,
		AccountOwnerID: accountOwnerID,
	}, nil
}

func (ucr *userAccountCacheRepository) Delete(ctx context.Context, userID string) error {
	key := ucr.buildUserAccountCacheKey(userID)

	if err := ucr.client.GetRDB().Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete user account cache error: %w", err)
	}
	return nil
}

func (ucr *userAccountCacheRepository) Exists(ctx context.Context, userID string) (bool, error) {
	key := ucr.buildUserAccountCacheKey(userID)

	count, err := ucr.client.GetRDB().Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists user account cache error: %w", err)
	}

	return count > 0, nil
}
