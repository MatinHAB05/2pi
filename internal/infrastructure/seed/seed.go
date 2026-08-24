package seed

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MatinHAB05/2pi/config"
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
	userRepo    repository_contract.UserRepository
	rbacRepo    repository_contract.RBACRepository
	accountRepo repository_contract.TargetAccountRepository
	cfg         *config.AdminConfig
}

func NewRBACSeeder(
	userRepo repository_contract.UserRepository,
	rbacRepo repository_contract.RBACRepository,
	accountRepo repository_contract.TargetAccountRepository,
	cfg *config.AdminConfig,
) *RBACSeeder {
	return &RBACSeeder{
		userRepo:    userRepo,
		rbacRepo:    rbacRepo,
		accountRepo: accountRepo,
		cfg:         cfg,
	}
}

func (s *RBACSeeder) SeedPermissions(ctx context.Context) error {
	for _, perm := range defaultPermissions {
		permissionAction := fmt.Sprintf("%s:%s", perm.Resource, perm.Action)

		_, err := s.rbacRepo.AddPermissionForRole(ctx, string(perm.Role), permissionAction)
		if err != nil {
			return fmt.Errorf("failed to seed permission [%s -> %s]: %w", perm.Role, permissionAction, err)
		}
	}
	return nil
}

func (s *RBACSeeder) SeedAdminUser(ctx context.Context) (int64, int64, error) {
	// 1. Create the Admin User entity itself
	adminUser := &entity.User{
		BaseEntity: entity.BaseEntity{ID: s.cfg.UserID},
		FirstName:  s.cfg.FirstName,
		LastName:   s.cfg.LastName,
		Username:   s.cfg.Username,
		UserLang:   entity.Lang(s.cfg.Language),
	}

	err := s.userRepo.Create(ctx, adminUser)
	if err != nil {
		return -1, -1, fmt.Errorf("failed to create admin user: %w", err)
	}

	adminUserID := adminUser.ID

	// 2. Create the target account owned by the created admin user
	adminAccount := &entity.TargetAccount{
		OwnerUserID: adminUserID,
		DayDuration: 97,
		Period:      97,
		Description: "this is admin",
		Enable:      false,
	}

	err = s.accountRepo.Create(ctx, adminAccount)
	if err != nil {
		return -1, -1, fmt.Errorf("failed to create admin target account: %w", err)
	}

	targetID := adminAccount.ID

	// 3. Assign role mapping in RBAC repository
	adminUserIDStr := strconv.FormatInt(adminUserID, 10)
	targetIDStr := strconv.FormatInt(targetID, 10)

	_, err = s.rbacRepo.AddUserRoleForTargetAccount(ctx, adminUserIDStr, targetIDStr, string(entity.RoleAdmin))
	if err != nil {
		return -1, -1, fmt.Errorf("failed to assign admin role to user %d with target %d: %w", adminUserID, targetID, err)
	}

	return adminUserID, targetID, nil
}

func (s *RBACSeeder) Execute(ctx context.Context) (int64, int64, error) {
	if err := s.SeedPermissions(ctx); err != nil {
		return -1, -1, err
	}

	adminUserID, targetID, err := s.SeedAdminUser(ctx)
	if err != nil {
		return -1, -1, err
	}

	return adminUserID, targetID, nil
}
