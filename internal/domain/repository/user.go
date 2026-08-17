package repository_contract

import (
	"context"

	"github.com/MatinHAB05/2pi/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	GetWithTargetAccounts(ctx context.Context, id uint) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error
	Exists(ctx context.Context, id uint) (bool, error)
}
