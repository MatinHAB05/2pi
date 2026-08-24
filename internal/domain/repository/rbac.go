package repository_contract

import "context"

type RBACRepository interface {
	EnforceForTargetAccount(ctx context.Context, userID string, targetAccountID string, action string) (bool, error)

	// User-Role Mapping (Grouping Policies - 'g: subject, resource, role')
	AddUserRoleForTargetAccount(ctx context.Context, userID string, targetAccountID string, role string) (bool, error)
	RemoveUserRoleForTargetAccount(ctx context.Context, userID string, targetAccountID string, role string) (bool, error)
	GetUserRolesForTargetAccount(ctx context.Context, userID string, targetAccountID string) ([]UserAccountRoleDTO, error)
	GetUsersForTargetAccount(ctx context.Context, targetAccountID string) ([]UserAccountRoleDTO, error)
	GetTargetAccountsForUser(ctx context.Context, userID string) ([]UserAccountRoleDTO, error)
	GetTargetAccountsForUsers(ctx context.Context, userID []string) ([]UserAccountRoleDTO, error)

	RemoveAllRolesForTargetAccount(ctx context.Context, targetAccountID string) (bool, error)
	RemoveAllRolesForUser(ctx context.Context, userID string) (bool, error)

	AddPermissionForRole(ctx context.Context, role string, action string) (bool, error)
	RemovePermissionForRole(ctx context.Context, role string, action string) (bool, error)
	GetPermissionsForRole(ctx context.Context, role string) ([]RolePermissionDTO, error)
}

type UserAccountRoleDTO struct {
	UserID          string
	TargetAccountID string
	Role            string
}

type RolePermissionDTO struct {
	Role   string
	Action string
}
