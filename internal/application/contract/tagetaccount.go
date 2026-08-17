package service_contract

import (
	"context"

	"github.com/MatinHAB05/2pi/internal/domain/entity"
)

type TargetAccountService interface {
	Create(ctx context.Context, req CreateTargetAccountRequest) (*TargetAccountResponse, error)
	GetByID(ctx context.Context, id uint) (*TargetAccountResponse, error)
	GetByUsername(ctx context.Context, username string) (*TargetAccountResponse, error)
	GetByOwnerID(ctx context.Context, ownerUserID int64, limit, offset int) ([]TargetAccountResponse, error)
	Update(ctx context.Context, req UpdateTargetAccountRequest) (*TargetAccountResponse, error)
	Delete(ctx context.Context, id uint) error
	CountByOwnerID(ctx context.Context, ownerUserID int64) (int64, error)
	DeleteByIDAndOwnerID(ctx context.Context, id uint, ownerUserID int64) error
	GetByIDAndOwnerID(ctx context.Context, id uint, ownerUserID int64) (*TargetAccountResponse, error)
}

type CreateTargetAccountRequest struct {
	OwnerUserID int64  `json:"owner_user_id"`
	Username    string `json:"username"`
}

type UpdateTargetAccountRequest struct {
	ID          uint   `json:"id"`
	OwnerUserID int64  `json:"owner_user_id"`
	Username    string `json:"username"`
	DayDuration int
	Period      int
	UserLang    entity.Lang
	Description string
	Enable      bool
}

type TargetAccountResponse struct {
	ID          uint   `json:"id"`
	OwnerUserID int64  `json:"owner_user_id"`
	Username    string `json:"username"`
	DayDuration int
	Period      int
	UserLang    entity.Lang
	Description string
	Enable      bool
}

func ToTargetAccountResponse(ta *entity.TargetAccount) *TargetAccountResponse {
	if ta == nil {
		return nil
	}
	return &TargetAccountResponse{
		ID:          uint(ta.ID),
		OwnerUserID: ta.OwnerUserID,
		Username:    ta.Username,
		DayDuration: *ta.DayDuration,
		Period:      *ta.DayDuration,
		UserLang:    entity.LangEng,
		Description: *ta.Description,
		Enable:      ta.Enable,
	}
}

func ToTargetAccountSliceResponse(accounts []entity.TargetAccount) []TargetAccountResponse {
	if len(accounts) == 0 {
		return nil
	}
	res := make([]TargetAccountResponse, len(accounts))
	for i := range accounts {
		if ta := ToTargetAccountResponse(&accounts[i]); ta != nil {
			res[i] = *ta
		}
	}
	return res
}
