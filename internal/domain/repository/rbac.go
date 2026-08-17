package repository_contract

type RBACRepository interface {
	EnforceForTargetAccount(userID string, targetAccountID string, action string) (bool, error)

	// User-Role Mapping (Grouping Policies - 'g')
	AddUserRoleForTargetAccount(userID string, role string, targetAccountID string) (bool, error)
	RemoveUserRoleForTargetAccount(userID string, role string, targetAccountID string) (bool, error)
	GetUserRolesForTargetAccount(userID string, targetAccountID string) ([]string, error)
	GetUsersForTargetAccount(targetAccountID string) ([][]string, error)
	RemoveAllRolesForTargetAccount(targetAccountID string) (bool, error)
	RemoveAllRolesForUser(userID string) (bool, error)

	// Dynamic Role-Permission Management (Policies - 'p')
	AddPermissionForRole(role string, action string) (bool, error)
	RemovePermissionForRole(role string, action string) (bool, error)
	GetPermissionsForRole(role string) ([][]string, error)
}
