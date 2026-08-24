package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/internal/infrastructure/database"
	"github.com/redis/go-redis/v9"
)

type shareAccountOTPCacheRepository struct {
	client database.Cache
}

func NewShareAccountOTPCacheRepository(client database.Cache) repository_contract.ShareAccountOTPRepository {
	return &shareAccountOTPCacheRepository{
		client: client,
	}
}

func (r *shareAccountOTPCacheRepository) buildShareAccountOTPCacheKey(otp string) string {
	return fmt.Sprintf("account:share:otp:%s", otp)
}

func (r *shareAccountOTPCacheRepository) Set(ctx context.Context, otp string, cache *repository_contract.ShareAccountOTP, ttl time.Duration) error {
	key := r.buildShareAccountOTPCacheKey(otp)

	fields := map[string]interface{}{
		"role":            cache.Role,
		"base_account_id": cache.BaseAccountID,
		"base_user_id":    cache.BaseUserID,
	}

	pipe := r.client.GetRDB().Pipeline()
	pipe.HSet(ctx, key, fields)
	if ttl > 0 {
		pipe.Expire(ctx, key, ttl)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("redis set otp cache error: %w", err)
	}

	return nil
}

func (r *shareAccountOTPCacheRepository) Get(ctx context.Context, otp string) (*repository_contract.ShareAccountOTP, error) {
	key := r.buildShareAccountOTPCacheKey(otp)

	res, err := r.client.GetRDB().HGetAll(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("redis get otp cache error: %w", err)
	}

	if len(res) == 0 {
		return nil, nil
	}

	return &repository_contract.ShareAccountOTP{
		Role:          res["role"],
		BaseAccountID: res["base_account_id"],
		BaseUserID:    res["base_user_id"],
	}, nil
}

func (r *shareAccountOTPCacheRepository) Delete(ctx context.Context, otp string) (bool, error) {
	key := r.buildShareAccountOTPCacheKey(otp)

	deletedCount, err := r.client.GetRDB().Del(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis delete otp cache error: %w", err)
	}

	return deletedCount > 0, nil
}

func (r *shareAccountOTPCacheRepository) Exists(ctx context.Context, otp string) (*bool, error) {
	key := r.buildShareAccountOTPCacheKey(otp)

	count, err := r.client.GetRDB().Exists(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis exists otp cache error: %w", err)
	}
	fin := count > 0
	return &fin, nil
}
