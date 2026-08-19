package service_contract

import (
	"context"

	"github.com/MatinHAB05/2pi/internal/domain/entity"
)

type TargetAccountService interface {
	Create(ctx context.Context, tokenContext TokenContext, req CreateTargetAccountRequest) (*TargetAccountResponse, error)
	GetByID(ctx context.Context, tokenContext TokenContext, id int64) (*TargetAccountResponse, error)
	GetByOwnerID(ctx context.Context, tokenContext TokenContext, ownerUserID int64, limit, offset int) ([]TargetAccountResponse, error)
	Update(ctx context.Context, tokenContext TokenContext, req UpdateTargetAccountRequest) (*TargetAccountResponse, error)
	UpdateStatus(ctx context.Context, tokenContext TokenContext, id int64, status bool) error
	Delete(ctx context.Context, tokenContext TokenContext, id int64) error
	CountByOwnerID(ctx context.Context, tokenContext TokenContext, ownerUserID int64) (int64, error)
	DeleteByIDAndOwnerID(ctx context.Context, tokenContext TokenContext, id int64, ownerUserID int64) error
	GetByIDAndOwnerID(ctx context.Context, tokenContext TokenContext, id int64, ownerUserID int64) (*TargetAccountResponse, error)
}

type CreateTargetAccountRequest struct {
	OwnerUserID int64 `json:"owner_user_id"`
}

type UpdateTargetAccountRequest struct {
	ID          int64 `json:"id"`
	OwnerUserID int64 `json:"owner_user_id"`
	DayDuration int
	Period      int
	Description string
	Enable      bool
}

type TargetAccountResponse struct {
	ID          int64 `json:"id"`
	OwnerUserID int64 `json:"owner_user_id"`
	DayDuration int
	Period      int
	Description string
	Enable      bool
}

func ToTargetAccountResponse(ta *entity.TargetAccount) *TargetAccountResponse {
	if ta == nil {
		return nil
	}
	return &TargetAccountResponse{
		ID:          int64(ta.ID),
		OwnerUserID: ta.OwnerUserID,
		DayDuration: ta.DayDuration,
		Period:      ta.Period,
		Description: ta.Description,
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
