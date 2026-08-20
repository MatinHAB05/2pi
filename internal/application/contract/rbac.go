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
	GetUserRolesForTargetAccount(ctx context.Context, tokenContext TokenContext, userID int64, targetAccountID int64) ([]*UserAccountRoleModelResponse, error)
	GetUsersForTargetAccount(ctx context.Context, tokenContext TokenContext, targetAccountID int64) ([]*UserAccountRoleModelResponse, error)
	GetTargetAccountsForUser(ctx context.Context, tokenContext TokenContext, userID int64) ([]*UserAccountRoleModelResponse, error)
	GetTargetAccountsForUsers(ctx context.Context, tokenContext TokenContext, userID []int64) ([]*UserAccountRoleModelResponse, error)
	RemoveAllRolesForTargetAccount(ctx context.Context, tokenContext TokenContext, targetAccountID int64) (bool, error)
	RemoveAllRolesForUser(ctx context.Context, tokenContext TokenContext, userID int64) (bool, error)

	// Dynamic Role-Permission Management
	AddPermissionForRole(ctx context.Context, tokenContext TokenContext, role string, action string) (bool, error)
	RemovePermissionForRole(ctx context.Context, tokenContext TokenContext, role string, action string) (bool, error)
	GetPermissionsForRole(ctx context.Context, tokenContext TokenContext, role string) ([]*RolePermissionModelResponse, error)
}

type UserAccountRoleModelResponse struct {
	UserID          int64  `json:"user_id"`
	TargetAccountID int64  `json:"target_account_id"`
	Role            string `json:"role"`
}

type RolePermissionModelResponse struct {
	Role   string `json:"role"`
	Action string `json:"action"`
}

func ToUserAccountRoleModelResponse(dto *repository_contract.UserAccountRoleDTO) *UserAccountRoleModelResponse {
	if dto == nil {
		return &UserAccountRoleModelResponse{}
	}
	userID, _ := strconv.ParseInt(dto.UserID, 10, 64)
	targetAccountID, _ := strconv.ParseInt(dto.TargetAccountID, 10, 64)

	return &UserAccountRoleModelResponse{
		UserID:          userID,
		TargetAccountID: targetAccountID,
		Role:            dto.Role,
	}
}

func ToUserAccountRoleSliceModelResponse(dtos []repository_contract.UserAccountRoleDTO) []*UserAccountRoleModelResponse {
	res := make([]*UserAccountRoleModelResponse, len(dtos))

	if len(dtos) == 0 {
		return res
	}
	for i := range dtos {
		if r := ToUserAccountRoleModelResponse(&dtos[i]); r != nil {
			res[i] = r
		}
	}
	return res
}

func ToRolePermissionModelResponse(dto *repository_contract.RolePermissionDTO) *RolePermissionModelResponse {
	if dto == nil {
		return &RolePermissionModelResponse{}
	}
	return &RolePermissionModelResponse{
		Role:   dto.Role,
		Action: dto.Action,
	}
}

func ToRolePermissionSliceModelModelResponse(dtos []repository_contract.RolePermissionDTO) []*RolePermissionModelResponse {
	res := make([]*RolePermissionModelResponse, len(dtos))

	if len(dtos) == 0 {
		return res
	}
	for i := range dtos {
		if r := ToRolePermissionModelResponse(&dtos[i]); r != nil {
			res[i] = r
		}
	}
	return res
}

func ExtractUserIDs(dtos []*UserAccountRoleModelResponse) []int64 {
	res := make([]int64, len(dtos))

	if len(dtos) == 0 {
		return res
	}

	for i := range dtos {
		res[i] = dtos[i].UserID
	}

	return res
}

func ExtractAccountIDs(dtos []*UserAccountRoleModelResponse) []int64 {
	res := make([]int64, len(dtos))

	if len(dtos) == 0 {
		return res
	}

	for i := range dtos {
		res[i] = dtos[i].TargetAccountID
	}

	return res
}
