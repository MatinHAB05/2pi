package service_contract

import "context"

type RBACService interface {
	EnforceForTargetAccount(ctx context.Context, userID string, targetAccountID string, action string) (bool, error)

	// User-Role Mapping
	AddUserRoleForTargetAccount(ctx context.Context, userID string, role string, targetAccountID string) (bool, error)
	RemoveUserRoleForTargetAccount(ctx context.Context, userID string, role string, targetAccountID string) (bool, error)
	GetUserRolesForTargetAccount(ctx context.Context, userID string, targetAccountID string) ([]string, error)
	GetUsersForTargetAccount(ctx context.Context, targetAccountID string) ([][]string, error)
	RemoveAllRolesForTargetAccount(ctx context.Context, targetAccountID string) (bool, error)
	RemoveAllRolesForUser(ctx context.Context, userID string) (bool, error)

	// Dynamic Role-Permission Management
	AddPermissionForRole(ctx context.Context, role string, action string) (bool, error)
	RemovePermissionForRole(ctx context.Context, role string, action string) (bool, error)
	GetPermissionsForRole(ctx context.Context, role string) ([][]string, error)
}
