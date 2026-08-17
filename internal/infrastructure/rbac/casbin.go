package rbac

import (
	"log"
	"sync"

	"github.com/MatinHAB05/reminder/config"
	"github.com/MatinHAB05/reminder/internal/infrastructure/database"
	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

type RBAC interface {
	GetEnforcer() *casbin.Enforcer
	Enforce(rvals ...any) (bool, error)
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

type Casbin struct {
	enforcer *casbin.Enforcer
}

var (
	casbinOnce     sync.Once
	casbinInstance *Casbin
)

func NewCasbin(casbinConfig *config.Casbin, gormDB database.Database) *Casbin {
	casbinOnce.Do(func() {
		adapter, err := gormadapter.NewAdapterByDB(gormDB.GetDB())
		if err != nil {
			log.Fatalf("Failed to create Casbin adapter: %v", err)
		}

		m, err := model.NewModelFromFile(casbinConfig.ModelConfigFilePath)
		if err != nil {
			log.Fatalf("Failed to parse Casbin model configuration: %v", err)
		}

		enforcer, err := casbin.NewEnforcer(m, adapter)
		if err != nil {
			log.Fatalf("Failed to create Casbin enforcer: %v", err)
		}

		err = enforcer.LoadPolicy()
		if err != nil {
			log.Fatalf("Failed to load Casbin policies from DB: %v", err)
		}

		casbinInstance = &Casbin{enforcer: enforcer}
	})

	return casbinInstance
}

// GetEnforcer returns the underlying Casbin enforcer instance.
func (c *Casbin) GetEnforcer() *casbin.Enforcer {
	return c.enforcer
}

// Enforce evaluates raw arguments against the Casbin model.
func (c *Casbin) Enforce(rvals ...any) (bool, error) {
	return c.enforcer.Enforce(rvals...)
}

// EnforceForTargetAccount evaluates access permissions for a user on a specific target account.
func (c *Casbin) EnforceForTargetAccount(userID string, targetAccountID string, action string) (bool, error) {
	return c.enforcer.Enforce(userID, targetAccountID, action)
}

// AddUserRoleForTargetAccount assigns a role to a user for a target account ('g' rule).
func (c *Casbin) AddUserRoleForTargetAccount(userID string, role string, targetAccountID string) (bool, error) {
	return c.enforcer.AddGroupingPolicy(userID, role, targetAccountID)
}

// RemoveUserRoleForTargetAccount revokes a role from a user for a target account.
func (c *Casbin) RemoveUserRoleForTargetAccount(userID string, role string, targetAccountID string) (bool, error) {
	return c.enforcer.RemoveGroupingPolicy(userID, role, targetAccountID)
}

// GetUserRolesForTargetAccount retrieves all roles of a user on a specific target account.
func (c *Casbin) GetUserRolesForTargetAccount(userID string, targetAccountID string) []string {
	return c.enforcer.GetRolesForUserInDomain(userID, targetAccountID)
}

// GetUsersForTargetAccount retrieves all user-role assignments for a target account.
func (c *Casbin) GetUsersForTargetAccount(targetAccountID string) ([][]string, error) {
	return c.enforcer.GetFilteredGroupingPolicy(2, targetAccountID)
}

// RemoveAllRolesForTargetAccount purges all role assignments for a deleted target account.
func (c *Casbin) RemoveAllRolesForTargetAccount(targetAccountID string) (bool, error) {
	return c.enforcer.RemoveFilteredGroupingPolicy(2, targetAccountID)
}

// RemoveAllRolesForUser purges all role assignments associated with a user.
func (c *Casbin) RemoveAllRolesForUser(userID string) (bool, error) {
	return c.enforcer.RemoveFilteredGroupingPolicy(0, userID)
}

// AddPermissionForRole dynamically grants an action permission to a role ('p' rule).
func (c *Casbin) AddPermissionForRole(role string, action string) (bool, error) {
	return c.enforcer.AddPolicy(role, action)
}

// RemovePermissionForRole dynamically revokes an action permission from a role.
func (c *Casbin) RemovePermissionForRole(role string, action string) (bool, error) {
	return c.enforcer.RemovePolicy(role, action)
}

// GetPermissionsForRole retrieves all assigned action permissions for a specific role.
func (c *Casbin) GetPermissionsForRole(role string) ([][]string, error) {
	return c.enforcer.GetFilteredPolicy(0, role)
}
