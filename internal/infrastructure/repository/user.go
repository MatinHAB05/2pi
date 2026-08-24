package repository

import (
	"context"
	"errors"

	"github.com/MatinHAB05/2pi/internal/domain/entity"
	"github.com/MatinHAB05/2pi/internal/domain/exception"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/internal/infrastructure/database"
	"gorm.io/gorm"
)

type userRepository struct {
	db database.Database
}

func NewUserRepository(db database.Database) repository_contract.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	db := database.ExtractTrxOrDB(ctx, r.db)
	return db.GetDB().WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	db := database.ExtractTrxOrDB(ctx, r.db)
	err := db.GetDB().WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrUserNotFound
	}
	return &user, err
}

func (r *userRepository) GetWithTargetAccounts(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	db := database.ExtractTrxOrDB(ctx, r.db)
	err := db.GetDB().WithContext(ctx).
		Preload("TargetAccounts").
		First(&user, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrUserNotFound
	}
	return &user, err
}

func (r *userRepository) GetByIDs(ctx context.Context, ids []int64) ([]entity.User, error) {
	var users []entity.User
	if len(ids) == 0 {
		return users, nil
	}

	db := database.ExtractTrxOrDB(ctx, r.db)
	err := db.GetDB().WithContext(ctx).Where("id IN ?", ids).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, err
}

func (r *userRepository) GetByIDsWithTargetAccounts(ctx context.Context, ids []int64) ([]entity.User, error) {
	var users []entity.User
	if len(ids) == 0 {
		return users, nil
	}

	db := database.ExtractTrxOrDB(ctx, r.db)
	err := db.GetDB().WithContext(ctx).Where("id IN ?", ids).Find(&users).Preload("TargetAccounts").Error
	if err != nil {
		return nil, err
	}
	return users, err
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	db := database.ExtractTrxOrDB(ctx, r.db)
	return db.GetDB().WithContext(ctx).Where("id = ?", user.BaseEntity.ID).Updates(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	db := database.ExtractTrxOrDB(ctx, r.db)
	result := db.GetDB().WithContext(ctx).Delete(&entity.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exception.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) Exists(ctx context.Context, id int64) (bool, error) {
	var count int64
	db := database.ExtractTrxOrDB(ctx, r.db)
	err := db.GetDB().WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", id).
		Count(&count).
		Error
	return count > 0, err
}
