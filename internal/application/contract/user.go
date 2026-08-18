package service_contract

import (
	"context"

	"github.com/MatinHAB05/2pi/internal/domain/entity"
)

type UserService interface {
	Create(ctx context.Context, tokenContext TokenContext, req CreateUserRequest) (*UserResponse, error)
	GetByID(ctx context.Context, tokenContext TokenContext, id int64) (*UserResponse, error)
	GetWithTargetAccounts(ctx context.Context, tokenContext TokenContext, id int64) (*UserResponse, error)
	Update(ctx context.Context, tokenContext TokenContext, req UpdateUserRequest) (*UserResponse, error)
	Delete(ctx context.Context, tokenContext TokenContext, id int64) error
	Exists(ctx context.Context, tokenContext TokenContext, id int64) (bool, error)
}

type CreateUserRequest struct {
	ID int64 `json:"id"`
}

type UpdateUserRequest struct {
	ID int64 `json:"id"`
}

type UserResponse struct {
	ID             int64                   `json:"id"`
	TargetAccounts []TargetAccountResponse `json:"target_accounts,omitempty"`
}

func ToUserResponse(u *entity.User) *UserResponse {
	if u == nil {
		return nil
	}
	res := &UserResponse{
		ID: int64(u.ID),
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
