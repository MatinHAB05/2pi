package service_contract

import (
	"context"
	"strconv"

	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
)

type RBACService interface {
	EnforceForTargetAccount(ctx context.Context, tokenContext TokenContext, userID int64, targetAccountID int64, action string) (bool, error)

	// User-Role Mapping
	AddUserRoleForTargetAccount(ctx context.Context, tokenContext TokenContext, userID int64, targetAccountID int64, role string) (bool, error)
	RemoveUserRoleForTargetAccount(ctx context.Context, tokenContext TokenContext, userID int64, targetAccountID int64, role string) (bool, error)
	GetUserRolesForTargetAccount(ctx context.Context, tokenContext TokenContext, userID int64, targetAccountID int64) ([]*UserAccountRoleResponse, error)
	GetUsersForTargetAccount(ctx context.Context, tokenContext TokenContext, targetAccountID int64) ([]*UserAccountRoleResponse, error)
	GetTargetAccountsForUser(ctx context.Context, tokenContext TokenContext, userID int64) ([]*UserAccountRoleResponse, error)
	RemoveAllRolesForTargetAccount(ctx context.Context, tokenContext TokenContext, targetAccountID int64) (bool, error)
	RemoveAllRolesForUser(ctx context.Context, tokenContext TokenContext, userID int64) (bool, error)

	// Dynamic Role-Permission Management
	AddPermissionForRole(ctx context.Context, tokenContext TokenContext, role string, action string) (bool, error)
	RemovePermissionForRole(ctx context.Context, tokenContext TokenContext, role string, action string) (bool, error)
	GetPermissionsForRole(ctx context.Context, tokenContext TokenContext, role string) ([]*RolePermissionResponse, error)
}

type UserAccountRoleResponse struct {
	UserID          int64  `json:"user_id"`
	TargetAccountID int64  `json:"target_account_id"`
	Role            string `json:"role"`
}

type RolePermissionResponse struct {
	Role   string `json:"role"`
	Action string `json:"action"`
}

func ToUserAccountRoleResponse(dto *repository_contract.UserAccountRoleDTO) *UserAccountRoleResponse {
	if dto == nil {
		return &UserAccountRoleResponse{}
	}
	userID, _ := strconv.ParseInt(dto.UserID, 10, 64)
	targetAccountID, _ := strconv.ParseInt(dto.TargetAccountID, 10, 64)

	return &UserAccountRoleResponse{
		UserID:          userID,
		TargetAccountID: targetAccountID,
		Role:            dto.Role,
	}
}

func ToUserAccountRoleSliceResponse(dtos []repository_contract.UserAccountRoleDTO) []*UserAccountRoleResponse {
	res := make([]*UserAccountRoleResponse, len(dtos))

	if len(dtos) == 0 {
		return res
	}
	for i := range dtos {
		if r := ToUserAccountRoleResponse(&dtos[i]); r != nil {
			res[i] = r
		}
	}
	return res
}

func ToRolePermissionResponse(dto *repository_contract.RolePermissionDTO) *RolePermissionResponse {
	if dto == nil {
		return &RolePermissionResponse{}
	}
	return &RolePermissionResponse{
		Role:   dto.Role,
		Action: dto.Action,
	}
}

func ToRolePermissionSliceResponse(dtos []repository_contract.RolePermissionDTO) []*RolePermissionResponse {
	res := make([]*RolePermissionResponse, len(dtos))

	if len(dtos) == 0 {
		return res
	}
	for i := range dtos {
		if r := ToRolePermissionResponse(&dtos[i]); r != nil {
			res[i] = r
		}
	}
	return res
}

func ExtractUserIDs(dtos []*UserAccountRoleResponse) []int64 {
	res := make([]int64, len(dtos))

	if len(dtos) == 0 {
		return res
	}

	for i := range dtos {
		res[i] = dtos[i].UserID
	}

	return res
}

func ExtractAccountIDs(dtos []*UserAccountRoleResponse) []int64 {
	res := make([]int64, len(dtos))

	if len(dtos) == 0 {
		return res
	}

	for i := range dtos {
		res[i] = dtos[i].TargetAccountID
	}

	return res
}
