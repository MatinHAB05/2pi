package service_contract

import "context"

type RBACService interface {
	EnforceForTargetAccount(ctx context.Context, tokenContext TokenContext, userID string, targetAccountID string, action string) (bool, error)

	// User-Role Mapping
	AddUserRoleForTargetAccount(ctx context.Context, tokenContext TokenContext, userID string, role string, targetAccountID string) (bool, error)
	RemoveUserRoleForTargetAccount(ctx context.Context, tokenContext TokenContext, userID string, role string, targetAccountID string) (bool, error)
	GetUserRolesForTargetAccount(ctx context.Context, tokenContext TokenContext, userID string, targetAccountID string) ([]string, error)
	GetUsersForTargetAccount(ctx context.Context, tokenContext TokenContext, targetAccountID string) ([][]string, error)
	RemoveAllRolesForTargetAccount(ctx context.Context, tokenContext TokenContext, targetAccountID string) (bool, error)
	RemoveAllRolesForUser(ctx context.Context, tokenContext TokenContext, userID string) (bool, error)

	// Dynamic Role-Permission Management
	AddPermissionForRole(ctx context.Context, tokenContext TokenContext, role string, action string) (bool, error)
	RemovePermissionForRole(ctx context.Context, tokenContext TokenContext, role string, action string) (bool, error)
	GetPermissionsForRole(ctx context.Context, tokenContext TokenContext, role string) ([][]string, error)
}
