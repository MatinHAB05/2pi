package repository

import (
	"context"
	"fmt"

	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/internal/infrastructure/database"
	"github.com/MatinHAB05/2pi/internal/infrastructure/rbac"
)

const (
	// Grouping Policy ('g') layout: g = subject, resource, role
	gIndexUserID          = 0
	gIndexTargetAccountID = 1
	gIndexRole            = 2

	// Permission Policy ('p') layout: p = role, action
	pIndexRole   = 0
	pIndexAction = 1
)

func buildGPolicy(userID, targetAccountID, role string) []string {
	g := make([]string, 3)
	g[gIndexUserID] = userID
	g[gIndexTargetAccountID] = targetAccountID
	g[gIndexRole] = role
	return g
}

func buildPPolicy(role, action string) []string {
	p := make([]string, 2)
	p[pIndexRole] = role
	p[pIndexAction] = action
	return p
}

type rbacRepository struct {
	db       database.Database
	enforcer rbac.RBACEnforcer
}

func NewRBACRepository(db database.Database, enforcer rbac.RBACEnforcer) repository_contract.RBACRepository {
	return &rbacRepository{db: db, enforcer: enforcer}
}

func (r *rbacRepository) EnforceForTargetAccount(userID string, targetAccountID string, action string) (bool, error) {
	ok, err := r.enforcer.GetEnforcer().Enforce(userID, targetAccountID, action)
	if err != nil {
		return false, fmt.Errorf("casbin enforce failed: %w", err)
	}
	return ok, nil
}

func (r *rbacRepository) AddUserRoleForTargetAccount(userID string, targetAccountID string, role string) (bool, error) {
	policy := buildGPolicy(userID, targetAccountID, role)
	ok, err := r.enforcer.GetEnforcer().AddGroupingPolicy(policy)
	if err != nil {
		return false, fmt.Errorf("casbin add grouping policy failed: %w", err)
	}
	return ok, nil
}

func (r *rbacRepository) RemoveUserRoleForTargetAccount(userID string, targetAccountID string, role string) (bool, error) {
	policy := buildGPolicy(userID, targetAccountID, role)
	ok, err := r.enforcer.GetEnforcer().RemoveGroupingPolicy(policy)
	if err != nil {
		return false, fmt.Errorf("casbin remove grouping policy failed: %w", err)
	}
	return ok, nil
}

// GetUserRolesForTargetAccount retrieves all user-role DTOs for a specific target account.
func (r *rbacRepository) GetUserRolesForTargetAccount(userID string, targetAccountID string) ([]repository_contract.UserAccountRoleDTO, error) {
	policies, err := r.enforcer.GetEnforcer().GetFilteredGroupingPolicy(gIndexUserID, userID)
	if err != nil {
		return nil, fmt.Errorf("casbin get filtered grouping policy failed: %w", err)
	}

	dtos := make([]repository_contract.UserAccountRoleDTO, 0, len(policies))
	for _, policy := range policies {
		if len(policy) > gIndexRole && len(policy) > gIndexTargetAccountID {
			if policy[gIndexTargetAccountID] == targetAccountID {
				dtos = append(dtos, repository_contract.UserAccountRoleDTO{
					UserID:          policy[gIndexUserID],
					TargetAccountID: policy[gIndexTargetAccountID],
					Role:            policy[gIndexRole],
				})
			}
		}
	}
	return dtos, nil
}

func (r *rbacRepository) GetUsersForTargetAccount(targetAccountID string) ([]repository_contract.UserAccountRoleDTO, error) {
	policies, err := r.enforcer.GetEnforcer().GetFilteredGroupingPolicy(gIndexTargetAccountID, targetAccountID)
	if err != nil {
		return nil, fmt.Errorf("casbin get filtered grouping policy failed: %w", err)
	}

	dtos := make([]repository_contract.UserAccountRoleDTO, 0, len(policies))
	for _, policy := range policies {
		if len(policy) > gIndexRole {
			dtos = append(dtos, repository_contract.UserAccountRoleDTO{
				UserID:          policy[gIndexUserID],
				TargetAccountID: policy[gIndexTargetAccountID],
				Role:            policy[gIndexRole],
			})
		}
	}
	return dtos, nil
}

func (r *rbacRepository) RemoveAllRolesForTargetAccount(targetAccountID string) (bool, error) {
	ok, err := r.enforcer.GetEnforcer().RemoveFilteredGroupingPolicy(gIndexTargetAccountID, targetAccountID)
	if err != nil {
		return false, fmt.Errorf("casbin remove filtered grouping policy for target failed: %w", err)
	}
	return ok, nil
}

func (r *rbacRepository) RemoveAllRolesForUser(userID string) (bool, error) {
	ok, err := r.enforcer.GetEnforcer().RemoveFilteredGroupingPolicy(gIndexUserID, userID)
	if err != nil {
		return false, fmt.Errorf("casbin remove filtered grouping policy for user failed: %w", err)
	}
	return ok, nil
}

func (r *rbacRepository) AddPermissionForRole(role string, action string) (bool, error) {
	policy := buildPPolicy(role, action)
	ok, err := r.enforcer.GetEnforcer().AddPolicy(policy)
	if err != nil {
		return false, fmt.Errorf("casbin add policy failed: %w", err)
	}
	return ok, nil
}

func (r *rbacRepository) RemovePermissionForRole(role string, action string) (bool, error) {
	policy := buildPPolicy(role, action)
	ok, err := r.enforcer.GetEnforcer().RemovePolicy(policy)
	if err != nil {
		return false, fmt.Errorf("casbin remove policy failed: %w", err)
	}
	return ok, nil
}

func (r *rbacRepository) GetPermissionsForRole(role string) ([]repository_contract.RolePermissionDTO, error) {
	policies, err := r.enforcer.GetEnforcer().GetFilteredPolicy(pIndexRole, role)
	if err != nil {
		return nil, fmt.Errorf("casbin get filtered policy failed: %w", err)
	}

	dtos := make([]repository_contract.RolePermissionDTO, 0, len(policies))
	for _, policy := range policies {
		if len(policy) > pIndexAction {
			dtos = append(dtos, repository_contract.RolePermissionDTO{
				Role:   policy[pIndexRole],
				Action: policy[pIndexAction],
			})
		}
	}
	return dtos, nil
}

func (r *rbacRepository) GetTargetAccountsForUser(ctx context.Context, userID string) ([]repository_contract.UserAccountRoleDTO, error) {
	policies, err := r.enforcer.GetEnforcer().GetFilteredGroupingPolicy(gIndexUserID, userID)
	if err != nil {
		return nil, fmt.Errorf("casbin get filtered grouping policy failed: %w", err)
	}

	dtos := make([]repository_contract.UserAccountRoleDTO, 0, len(policies))
	for _, policy := range policies {
		if len(policy) > gIndexRole {
			dtos = append(dtos, repository_contract.UserAccountRoleDTO{
				UserID:          policy[gIndexUserID],
				TargetAccountID: policy[gIndexTargetAccountID],
				Role:            policy[gIndexRole],
			})
		}
	}
	return dtos, nil
}

//? bug that i fix it .
// the fieldValues is not like IN operator instead it is used for multiple filter base on different fields(where field1==val1 && field2==val2 ...)
//func (e *Enforcer) GetFilteredGroupingPolicy(fieldIndex int, fieldValues ...string) ([][]string, error) {

func (r *rbacRepository) GetTargetAccountsForUsers(ctx context.Context, userIDs []string) ([]repository_contract.UserAccountRoleDTO, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	// 1. Build a lookup set for fast IN filtering
	userSet := make(map[string]struct{}, len(userIDs))
	for _, id := range userIDs {
		userSet[id] = struct{}{}
	}

	// 2. Fetch all grouping policies in a single call
	policies, err := r.enforcer.GetEnforcer().GetGroupingPolicy()
	if err != nil {
		return nil, fmt.Errorf("casbin get grouping policy failed: %w", err)
	}

	// 3. Filter policies matching any ID in userSet
	var dtos []repository_contract.UserAccountRoleDTO
	for _, policy := range policies {
		if len(policy) <= gIndexRole {
			continue
		}

		userID := policy[gIndexUserID]
		if _, exists := userSet[userID]; exists {
			dtos = append(dtos, repository_contract.UserAccountRoleDTO{
				UserID:          userID,
				TargetAccountID: policy[gIndexTargetAccountID],
				Role:            policy[gIndexRole],
			})
		}
	}

	return dtos, nil
}
