package repository

//BUG : replace string instead of set!!!!!! zero get check! +++ dto === account id , owner id
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
	return fmt.Sprintf("user:%s:accounts", userId)
}

func (ucr *useraccountcahceRepository) CreateOrReplace(ctx context.Context, userId string, ttl time.Duration) error {
	key := ucr.buildUserAccountCacheKey(userId)

	err := ucr.Delete(ctx, userId)
	if err != nil {
		return err
	}
	intUserId, err := strconv.Atoi(userId)
	if err != nil {
		return fmt.Errorf("WTF:%w", err)
	}

	user, err := ucr.userRepo.GetWithTargetAccounts(ctx, int64(intUserId))
	if err != nil {
		return err
	}
	accIds := make([]int64, len(user.TargetAccounts))
	for i := 0; i < len(user.TargetAccounts); i++ {
		accIds = append(accIds, user.TargetAccounts[i].ID)
	}

	if err := ucr.client.GetRDB().SAdd(ctx, key, accIds).Err(); err != nil {
		return fmt.Errorf("redis save user-account error: %w", err)
	}
	return nil
}

func (ucr *useraccountcahceRepository) CreateOrReplaceGet(ctx context.Context, userId string, ttl time.Duration) ([]string, error) {
	if err := ucr.CreateOrReplace(ctx, userId, ttl); err != nil {
		return nil, err
	}

	accs, err := ucr.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	return accs, nil

}

func (ucr *useraccountcahceRepository) Get(ctx context.Context, userId string) ([]string, error) {
	key := ucr.buildUserAccountCacheKey(userId)

	accs, err := ucr.client.GetRDB().SMembers(ctx, key).Result()
	if err != nil || len(accs) == 0 {
		if errors.Is(err, redis.Nil) || len(accs) == 0 {
			return []string{}, nil
		}
		return nil, fmt.Errorf("redis get user permissions error: %w", err)
	}

	return accs, nil
}
func (ucr *useraccountcahceRepository) Delete(ctx context.Context, userId string) error {
	key := ucr.buildUserAccountCacheKey(userId)

	if err := ucr.client.GetRDB().Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete user-account error: %w", err)
	}
	return nil
}
func (ucr *useraccountcahceRepository) Exists(ctx context.Context, userId string) (*bool, error) {
	key := ucr.buildUserAccountCacheKey(userId)

	count, err := ucr.client.GetRDB().Exists(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis exists user-account error: %w", err)
	}

	exists := count > 0
	return &exists, nil
}

func (ucr *useraccountcahceRepository) GetSync(ctx context.Context, userId string, ttl time.Duration) ([]string, error) {
	var accs []string
	ex, err := ucr.Exists(ctx, userId)
	if err != nil {
		return accs, err
	}
	if !*ex {
		accs, err = ucr.CreateOrReplaceGet(ctx, userId, ttl)
		return accs, nil
	}

	return accs, nil
}
