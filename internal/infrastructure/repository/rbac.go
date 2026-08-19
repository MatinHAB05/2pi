package repository

import (
	"context"
	"fmt"

	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/internal/infrastructure/database"
	"github.com/MatinHAB05/2pi/internal/infrastructure/rbac"
)

type rbacRepository struct {
	db       database.Database
	enforcer rbac.RBACEnforcer
}

func NewRBACRepository(db database.Database, enforcer rbac.RBACEnforcer) repository_contract.RBACRepository {
	return &rbacRepository{db: db, enforcer: enforcer}
}

// EnforceForTargetAccount evaluates access permissions for a user on a specific target account.
func (r *rbacRepository) EnforceForTargetAccount(userID string, targetAccountID string, action string) (bool, error) {
	ok, err := r.enforcer.GetEnforcer().Enforce(userID, targetAccountID, action)
	if err != nil {
		return false, fmt.Errorf("casbin enforce failed: %w", err)
	}
	return ok, nil
}

// AddUserRoleForTargetAccount assigns a role to a user for a target account ('g' rule: subject, resource, role).
func (r *rbacRepository) AddUserRoleForTargetAccount(userID string, targetAccountID string, role string) (bool, error) {
	ok, err := r.enforcer.GetEnforcer().AddGroupingPolicy(userID, targetAccountID, role)
	if err != nil {
		return false, fmt.Errorf("casbin add grouping policy failed: %w", err)
	}
	return ok, nil
}

// RemoveUserRoleForTargetAccount revokes a role from a user for a target account ('g' rule: subject, resource, role).
func (r *rbacRepository) RemoveUserRoleForTargetAccount(userID string, targetAccountID string, role string) (bool, error) {
	ok, err := r.enforcer.GetEnforcer().RemoveGroupingPolicy(userID, targetAccountID, role)
	if err != nil {
		return false, fmt.Errorf("casbin remove grouping policy failed: %w", err)
	}
	return ok, nil
}

// GetUserRolesForTargetAccount retrieves all roles of a user on a specific target account.
func (r *rbacRepository) GetUserRolesForTargetAccount(userID string, targetAccountID string) ([]string, error) {
	roles := r.enforcer.GetEnforcer().GetRolesForUserInDomain(userID, targetAccountID)

	return roles, nil
}

// GetUsersForTargetAccount retrieves all user-role assignments for a target account.
func (r *rbacRepository) GetUsersForTargetAccount(targetAccountID string) ([][]string, error) {
	policies, err := r.enforcer.GetEnforcer().GetFilteredGroupingPolicy(2, targetAccountID)
	if err != nil {
		return nil, fmt.Errorf("casbin get filtered grouping policy failed: %w", err)
	}
	return policies, nil
}

// RemoveAllRolesForTargetAccount purges all role assignments for a deleted target account.
func (r *rbacRepository) RemoveAllRolesForTargetAccount(targetAccountID string) (bool, error) {
	ok, err := r.enforcer.GetEnforcer().RemoveFilteredGroupingPolicy(2, targetAccountID)
	if err != nil {
		return false, fmt.Errorf("casbin remove filtered grouping policy for target failed: %w", err)
	}
	return ok, nil
}

// RemoveAllRolesForUser purges all role assignments associated with a user.
func (r *rbacRepository) RemoveAllRolesForUser(userID string) (bool, error) {
	ok, err := r.enforcer.GetEnforcer().RemoveFilteredGroupingPolicy(0, userID)
	if err != nil {
		return false, fmt.Errorf("casbin remove filtered grouping policy for user failed: %w", err)
	}
	return ok, nil
}

// AddPermissionForRole dynamically grants an action permission to a role ('p' rule).
func (r *rbacRepository) AddPermissionForRole(role string, action string) (bool, error) {
	ok, err := r.enforcer.GetEnforcer().AddPolicy(role, action)
	if err != nil {
		return false, fmt.Errorf("casbin add policy failed: %w", err)
	}
	return ok, nil
}

// RemovePermissionForRole dynamically revokes an action permission from a role.
func (r *rbacRepository) RemovePermissionForRole(role string, action string) (bool, error) {
	ok, err := r.enforcer.GetEnforcer().RemovePolicy(role, action)
	if err != nil {
		return false, fmt.Errorf("casbin remove policy failed: %w", err)
	}
	return ok, nil
}

// GetPermissionsForRole retrieves all assigned action permissions for a specific role.
func (r *rbacRepository) GetPermissionsForRole(role string) ([][]string, error) {
	policies, err := r.enforcer.GetEnforcer().GetFilteredPolicy(0, role)
	if err != nil {
		return nil, fmt.Errorf("casbin get filtered policy failed: %w", err)
	}
	return policies, nil
}

func (r *rbacRepository) GetTargetAccountsForUser(ctx context.Context, userID string) ([][]string, error) {
	roles, err := r.enforcer.GetEnforcer().GetFilteredGroupingPolicy(0, userID)
	if err != nil {
		return nil, fmt.Errorf("casbin get filtered policy failed: %w", err)
	}
	return roles, nil
}
