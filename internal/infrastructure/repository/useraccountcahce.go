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

type useraccountcahceRepository struct {
	client   database.Cache
	userRepo repository_contract.UserRepository
}

func NewUserAccountCacheRepository(
	client database.Cache,
	userRepo repository_contract.UserRepository,
) repository_contract.UserAccountCacheRepository {
	return &useraccountcahceRepository{
		client:   client,
		userRepo: userRepo,
	}
}

func (ucr *useraccountcahceRepository) buildUserAccountCacheKey(userId string) string {
	return fmt.Sprintf("user:%s:account", userId)
}

func (ucr *useraccountcahceRepository) CreateOrReplace(ctx context.Context, userId string, ttl time.Duration) error {
	key := ucr.buildUserAccountCacheKey(userId)

	// Clear existing cache key
	if err := ucr.Delete(ctx, userId); err != nil {
		return err
	}

	intUserId, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user_id format: %w", err)
	}

	user, err := ucr.userRepo.GetWithTargetAccounts(ctx, intUserId)
	if err != nil {
		return err
	}

	if len(user.TargetAccounts) == 0 {
		return nil
	}

	// Single target account mapping
	targetAcc := user.TargetAccounts[0]
	fields := map[string]interface{}{
		"account_id":       targetAcc.ID,
		"account_owner_id": targetAcc.OwnerUserID,
	}

	// Store fields as Redis Hash
	if err := ucr.client.GetRDB().HSet(ctx, key, fields).Err(); err != nil {
		return fmt.Errorf("redis save user-account hash error: %w", err)
	}

	// Apply TTL if specified
	if ttl > 0 {
		if err := ucr.client.GetRDB().Expire(ctx, key, ttl).Err(); err != nil {
			return fmt.Errorf("redis set expire error: %w", err)
		}
	}

	return nil
}

func (ucr *useraccountcahceRepository) CreateOrReplaceGet(ctx context.Context, userId string, ttl time.Duration) (*repository_contract.UserAccountCache, error) {
	if err := ucr.CreateOrReplace(ctx, userId, ttl); err != nil {
		return nil, err
	}

	return ucr.Get(ctx, userId)
}

func (ucr *useraccountcahceRepository) Get(ctx context.Context, userId string) (*repository_contract.UserAccountCache, error) {
	key := ucr.buildUserAccountCacheKey(userId)

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

	AccountOwnerID, err := strconv.ParseInt(res["account_owner_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid account_owner_id in redis: %w", err)
	}

	return &repository_contract.UserAccountCache{
		AccountID:      accountID,
		AccountOwnerID: AccountOwnerID,
	}, nil
}

func (ucr *useraccountcahceRepository) Delete(ctx context.Context, userId string) error {
	key := ucr.buildUserAccountCacheKey(userId)

	if err := ucr.client.GetRDB().Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete user-account error: %w", err)
	}
	return nil
}

func (ucr *useraccountcahceRepository) Exists(ctx context.Context, userId string) (bool, error) {
	key := ucr.buildUserAccountCacheKey(userId)

	count, err := ucr.client.GetRDB().Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists user-account error: %w", err)
	}

	return count > 0, nil
}

func (ucr *useraccountcahceRepository) GetSync(ctx context.Context, userId string, ttl time.Duration) (*repository_contract.UserAccountCache, error) {
	exists, err := ucr.Exists(ctx, userId)
	if err != nil {
		return nil, err
	}

	if !exists {
		return ucr.CreateOrReplaceGet(ctx, userId, ttl)
	}

	return ucr.Get(ctx, userId)
}
