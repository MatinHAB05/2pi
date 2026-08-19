package repository_contract

import "context"

type RBACRepository interface {
	EnforceForTargetAccount(userID string, targetAccountID string, action string) (bool, error)

	// User-Role Mapping (Grouping Policies - 'g: subject, resource, role')
	AddUserRoleForTargetAccount(userID string, targetAccountID string, role string) (bool, error)
	RemoveUserRoleForTargetAccount(userID string, targetAccountID string, role string) (bool, error)
	GetUserRolesForTargetAccount(userID string, targetAccountID string) ([]string, error)
	GetUsersForTargetAccount(targetAccountID string) ([][]string, error)
	GetTargetAccountsForUser(ctx context.Context, userID string) ([][]string, error)

	RemoveAllRolesForTargetAccount(targetAccountID string) (bool, error)
	RemoveAllRolesForUser(userID string) (bool, error)

	// Dynamic Role-Permission Management (Policies - 'p: role, action')
	AddPermissionForRole(role string, action string) (bool, error)
	RemovePermissionForRole(role string, action string) (bool, error)
	GetPermissionsForRole(role string) ([][]string, error)
}
