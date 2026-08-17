package repository_contract

import (
	"context"

	"github.com/MatinHAB05/reminder/internal/domain/entity"
)

type TargetAccountRepository interface {
	Create(ctx context.Context, target *entity.TargetAccount) error
	GetByID(ctx context.Context, id uint) (*entity.TargetAccount, error)
	GetByUsername(ctx context.Context, username string) (*entity.TargetAccount, error)
	GetByOwnerID(ctx context.Context, ownerUserID int64, limit, offset int) ([]entity.TargetAccount, error)
	Update(ctx context.Context, target *entity.TargetAccount) error
	Delete(ctx context.Context, id uint) error
	CountByOwnerID(ctx context.Context, ownerUserID int64) (int64, error)
	DeleteByIDAndOwnerID(ctx context.Context, id uint, ownerUserID int64) error
	GetByIDAndOwnerID(ctx context.Context, id uint, ownerUserID int64) (*entity.TargetAccount, error)
}
