package service_contract

import (
	"context"

	"github.com/MatinHAB05/2pi/internal/domain/entity"
)

type UserService interface {
	Create(ctx context.Context, tokenContext TokenContext, req CreateUserRequest) (*UserResponse, error)
	GetByID(ctx context.Context, tokenContext TokenContext, id int64) (*UserResponse, error)
	GetWithTargetAccounts(ctx context.Context, tokenContext TokenContext, id int64) (*UserResponse, error)
	GetByIDs(ctx context.Context, tokenContext TokenContext, ids []int64) ([]*UserResponse, error)
	GetByIDsWithTargetAccounts(ctx context.Context, tokenContext TokenContext, ids []int64) ([]*UserResponse, error)
	Update(ctx context.Context, tokenContext TokenContext, req UpdateUserRequest) (*UserResponse, error)
	Delete(ctx context.Context, tokenContext TokenContext, id int64) error
	Exists(ctx context.Context, tokenContext TokenContext, id int64) (bool, error)

	GetUserAccountsRolesByID(ctx context.Context, tokenContext TokenContext, id int64) (*UserAccountsRoleResponse, error)
	GetUserAccountsRolesByIDs(ctx context.Context, tokenContext TokenContext, ids []int64) ([]*UserAccountsRoleResponse, error)
}

type CreateUserRequest struct {
	ID        int64 `json:"id"`
	FirstName string
	LastName  string
	Username  string
}

type UpdateUserRequest struct {
	ID        int64       `json:"id"`
	Lang      entity.Lang `json:"lang"`
	FirstName string
	LastName  string
	Username  string
}

type UserResponse struct {
	ID             int64       `json:"id"`
	Lang           entity.Lang `json:"lang"`
	FirstName      string
	LastName       string
	Username       string
	TargetAccounts []*TargetAccountResponse `json:"target_accounts,omitempty"`
}
type UserAccountsRoleResponse struct {
	*UserResponse
	AccountsRoles []*AccountsRoleResponse `json:",omitempty"`
}

type AccountsRoleResponse struct {
	AccountID *TargetAccountResponse
	Roles     []string
}

func ToUserResponse(u *entity.User) *UserResponse {
	if u == nil {
		return &UserResponse{}
	}
	res := &UserResponse{
		ID:        int64(u.ID),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Lang:      u.UserLang,
	}
	if len(u.TargetAccounts) > 0 {
		res.TargetAccounts = ToTargetAccountSliceResponse(u.TargetAccounts)
	}

	return res
}

func ToUserSliceResponse(users []entity.User) []*UserResponse {
	res := make([]*UserResponse, len(users))

	if len(users) == 0 {
		return res
	}
	for i := range users {
		if u := ToUserResponse(&users[i]); u != nil {
			res[i] = u
		}
	}
	return res
}
