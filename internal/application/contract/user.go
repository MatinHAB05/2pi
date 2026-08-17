package service_contract

import (
	"context"

	"github.com/MatinHAB05/2pi/internal/domain/entity"
)

type UserService interface {
	Create(ctx context.Context, req CreateUserRequest) (*UserResponse, error)
	GetByID(ctx context.Context, id uint) (*UserResponse, error)
	GetWithTargetAccounts(ctx context.Context, id uint) (*UserResponse, error)
	Update(ctx context.Context, req UpdateUserRequest) (*UserResponse, error)
	Delete(ctx context.Context, id uint) error
	Exists(ctx context.Context, id uint) (bool, error)
}

type CreateUserRequest struct {
	ID uint `json:"id"`
}

type UpdateUserRequest struct {
	ID uint `json:"id"`
}

type UserResponse struct {
	ID             uint                    `json:"id"`
	TargetAccounts []TargetAccountResponse `json:"target_accounts,omitempty"`
}

func ToUserResponse(u *entity.User) *UserResponse {
	if u == nil {
		return nil
	}
	res := &UserResponse{
		ID: uint(u.ID),
	}
	if len(u.TargetAccounts) > 0 {
		res.TargetAccounts = ToTargetAccountSliceResponse(u.TargetAccounts)
	}
	return res
}

func ToUserSliceResponse(users []entity.User) []UserResponse {
	if len(users) == 0 {
		return nil
	}
	res := make([]UserResponse, len(users))
	for i := range users {
		if u := ToUserResponse(&users[i]); u != nil {
			res[i] = *u
		}
	}
	return res
}
