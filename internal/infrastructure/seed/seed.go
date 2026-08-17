package seed

import (
	"context"
	"fmt"

	"github.com/MatinHAB05/2pi/internal/domain/entity"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
)

type permissionRule struct {
	Role     entity.Role
	Resource entity.Resource
	Action   entity.Action
}

var defaultPermissions = []permissionRule{
	// Admin:
	{Role: entity.RoleAdmin, Resource: entity.ResourceAccount, Action: entity.ActionRead},
	{Role: entity.RoleAdmin, Resource: entity.ResourceAccount, Action: entity.ActionWrite},
	{Role: entity.RoleAdmin, Resource: entity.ResourceAccount, Action: entity.ActionDelete},
	{Role: entity.RoleAdmin, Resource: entity.ResourcePeriod, Action: entity.ActionRead},
	{Role: entity.RoleAdmin, Resource: entity.ResourcePeriod, Action: entity.ActionWrite},
	{Role: entity.RoleAdmin, Resource: entity.ResourcePeriod, Action: entity.ActionDelete},
	{Role: entity.RoleAdmin, Resource: entity.ResourceAccess, Action: entity.ActionRead},
	{Role: entity.RoleAdmin, Resource: entity.ResourceAccess, Action: entity.ActionWrite},
	{Role: entity.RoleAdmin, Resource: entity.ResourceAccess, Action: entity.ActionDelete},

	// Owner:
	{Role: entity.RoleOwner, Resource: entity.ResourceAccount, Action: entity.ActionRead},
	{Role: entity.RoleOwner, Resource: entity.ResourceAccount, Action: entity.ActionWrite},
	{Role: entity.RoleOwner, Resource: entity.ResourceAccount, Action: entity.ActionDelete},
	{Role: entity.RoleOwner, Resource: entity.ResourcePeriod, Action: entity.ActionRead},
	{Role: entity.RoleOwner, Resource: entity.ResourcePeriod, Action: entity.ActionWrite},
	{Role: entity.RoleOwner, Resource: entity.ResourcePeriod, Action: entity.ActionDelete},
	{Role: entity.RoleOwner, Resource: entity.ResourceAccess, Action: entity.ActionRead},
	{Role: entity.RoleOwner, Resource: entity.ResourceAccess, Action: entity.ActionWrite},
	{Role: entity.RoleOwner, Resource: entity.ResourceAccess, Action: entity.ActionDelete},

	// Editor:
	{Role: entity.RoleEditor, Resource: entity.ResourceAccount, Action: entity.ActionRead},
	{Role: entity.RoleEditor, Resource: entity.ResourceAccount, Action: entity.ActionWrite},
	{Role: entity.RoleEditor, Resource: entity.ResourcePeriod, Action: entity.ActionRead},
	{Role: entity.RoleEditor, Resource: entity.ResourcePeriod, Action: entity.ActionWrite},
	{Role: entity.RoleEditor, Resource: entity.ResourceAccess, Action: entity.ActionRead},

	// Viewer:
	{Role: entity.RoleViewer, Resource: entity.ResourceAccount, Action: entity.ActionRead},
	{Role: entity.RoleViewer, Resource: entity.ResourcePeriod, Action: entity.ActionRead},
	{Role: entity.RoleViewer, Resource: entity.ResourceAccess, Action: entity.ActionRead},
}

type RBACSeeder struct {
	rbacRepo repository_contract.RBACRepository
}

func NewRBACSeeder(rbacRepo repository_contract.RBACRepository) *RBACSeeder {
	return &RBACSeeder{rbacRepo: rbacRepo}
}

func (s *RBACSeeder) SeedPermissions(ctx context.Context) error {
	for _, perm := range defaultPermissions {
		permissionAction := fmt.Sprintf("%s:%s", perm.Resource, perm.Action)

		_, err := s.rbacRepo.AddPermissionForRole(string(perm.Role), permissionAction)
		if err != nil {
			return fmt.Errorf("failed to seed permission [%s -> %s]: %w", perm.Role, permissionAction, err)
		}
	}
	return nil
}

func (s *RBACSeeder) SeedAdminUser(ctx context.Context, adminUserID string, globalTargetID string) error {
	_, err := s.rbacRepo.AddUserRoleForTargetAccount(adminUserID, string(entity.RoleAdmin), globalTargetID)
	if err != nil {
		return fmt.Errorf("failed to assign admin role to user %s: %w", adminUserID, err)
	}
	return nil
}

func (s *RBACSeeder) Execute(ctx context.Context, adminUserID string, globalTargetID string) error {
	if err := s.SeedPermissions(ctx); err != nil {
		return err
	}
	if err := s.SeedAdminUser(ctx, adminUserID, globalTargetID); err != nil {
		return err
	}
	return nil
}
