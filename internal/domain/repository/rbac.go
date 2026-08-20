package repository_contract

import "context"

type RBACRepository interface {
	EnforceForTargetAccount(userID string, targetAccountID string, action string) (bool, error)

	// User-Role Mapping (Grouping Policies - 'g: subject, resource, role')
	AddUserRoleForTargetAccount(userID string, targetAccountID string, role string) (bool, error)
	RemoveUserRoleForTargetAccount(userID string, targetAccountID string, role string) (bool, error)
	GetUserRolesForTargetAccount(userID string, targetAccountID string) ([]UserAccountRoleDTO, error)
	GetUsersForTargetAccount(targetAccountID string) ([]UserAccountRoleDTO, error)
	GetTargetAccountsForUser(ctx context.Context, userID string) ([]UserAccountRoleDTO, error)
	GetTargetAccountsForUsers(ctx context.Context, userID []string) ([]UserAccountRoleDTO, error)

	RemoveAllRolesForTargetAccount(targetAccountID string) (bool, error)
	RemoveAllRolesForUser(userID string) (bool, error)

	AddPermissionForRole(role string, action string) (bool, error)
	RemovePermissionForRole(role string, action string) (bool, error)
	GetPermissionsForRole(role string) ([]RolePermissionDTO, error)
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
