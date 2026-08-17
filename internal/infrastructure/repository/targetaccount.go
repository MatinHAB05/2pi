package repository

import (
	"context"
	"errors"

	"github.com/MatinHAB05/2pi/internal/domain/entity"
	"github.com/MatinHAB05/2pi/internal/domain/exception"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"gorm.io/gorm"
)

type targetAccountRepository struct {
	db *gorm.DB
}

func NewTargetAccountRepository(db *gorm.DB) repository_contract.TargetAccountRepository {
	return &targetAccountRepository{db: db}
}

func (r *targetAccountRepository) Create(ctx context.Context, target *entity.TargetAccount) error {
	return r.db.WithContext(ctx).Create(target).Error
}

func (r *targetAccountRepository) GetByID(ctx context.Context, id uint) (*entity.TargetAccount, error) {
	var target entity.TargetAccount
	err := r.db.WithContext(ctx).First(&target, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrTargetAccountNotFound
	}
	return &target, err
}

func (r *targetAccountRepository) GetByUsername(ctx context.Context, username string) (*entity.TargetAccount, error) {
	var target entity.TargetAccount
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&target).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrTargetAccountNotFound
	}
	return &target, err
}

func (r *targetAccountRepository) GetByOwnerID(ctx context.Context, ownerUserID int64, limit, offset int) ([]entity.TargetAccount, error) {
	var targets []entity.TargetAccount
	query := r.db.WithContext(ctx).Where("owner_user_id = ?", ownerUserID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&targets).Error
	return targets, err
}

func (r *targetAccountRepository) GetByIDAndOwnerID(ctx context.Context, id uint, ownerUserID int64) (*entity.TargetAccount, error) {
	var target entity.TargetAccount
	err := r.db.WithContext(ctx).
		Where("id = ? AND owner_user_id = ?", id, ownerUserID).
		First(&target).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrTargetAccountNotFound
	}
	return &target, err
}

func (r *targetAccountRepository) Update(ctx context.Context, target *entity.TargetAccount) error {
	return r.db.WithContext(ctx).Save(target).Error
}

func (r *targetAccountRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&entity.TargetAccount{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exception.ErrTargetAccountNotFound
	}
	return nil
}

func (r *targetAccountRepository) DeleteByIDAndOwnerID(ctx context.Context, id uint, ownerUserID int64) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND owner_user_id = ?", id, ownerUserID).
		Delete(&entity.TargetAccount{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exception.ErrTargetAccountNotFound
	}
	return nil
}

func (r *targetAccountRepository) CountByOwnerID(ctx context.Context, ownerUserID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.TargetAccount{}).
		Where("owner_user_id = ?", ownerUserID).
		Count(&count).
		Error
	return count, err
}
