package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MatinHAB05/2pi/internal/domain/entity"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/internal/infrastructure/database"
	"github.com/redis/go-redis/v9"
)

type userInfoCacheRepository struct {
	client database.Cache
}

func NewUserInfoCacheRepository(client database.Cache) repository_contract.UserInfoCacheRepository {
	return &userInfoCacheRepository{
		client: client,
	}
}

func (ucr *userInfoCacheRepository) buildUserInfoCacheKey(userID string) string {
	return fmt.Sprintf("user:%s:info", userID)
}

func (ucr *userInfoCacheRepository) Set(ctx context.Context, userID string, cache *repository_contract.UserInfoCache, ttl time.Duration) error {
	key := ucr.buildUserInfoCacheKey(userID)

	fields := map[string]interface{}{
		"lang": string(cache.Lang),
	}

	pipe := ucr.client.GetRDB().Pipeline()
	pipe.HSet(ctx, key, fields)
	if ttl > 0 {
		pipe.Expire(ctx, key, ttl)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("redis set user info cache error: %w", err)
	}

	return nil
}

func (ucr *userInfoCacheRepository) Get(ctx context.Context, userID string) (*repository_contract.UserInfoCache, error) {
	key := ucr.buildUserInfoCacheKey(userID)

	res, err := ucr.client.GetRDB().HGetAll(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("redis get user info error: %w", err)
	}

	if len(res) == 0 {
		return nil, nil
	}

	return &repository_contract.UserInfoCache{
		Lang: entity.Lang(res["lang"]),
	}, nil
}

func (ucr *userInfoCacheRepository) Delete(ctx context.Context, userID string) error {
	key := ucr.buildUserInfoCacheKey(userID)

	if err := ucr.client.GetRDB().Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete user info cache error: %w", err)
	}
	return nil
}

func (ucr *userInfoCacheRepository) Exists(ctx context.Context, userID string) (*bool, error) {
	key := ucr.buildUserInfoCacheKey(userID)

	count, err := ucr.client.GetRDB().Exists(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis exists user info cache error: %w", err)
	}
	res := count > 0
	return &res, nil
}
