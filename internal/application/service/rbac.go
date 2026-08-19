package service

import (
	"context"
	"fmt"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/entity"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/pkg/logger"
)

// TODO : seed = action-role
type rbacService struct {
	repo   repository_contract.RBACRepository
	logger logger.Logger
}

func NewRBACService(repo repository_contract.RBACRepository, log logger.Logger) service_contract.RBACService {
	return &rbacService{
		repo:   repo,
		logger: log,
	}
}

func (s *rbacService) isValidRole(role string) bool {
	return entity.RoleSet[role]
}

func (s *rbacService) isValidAction(action string) bool {
	return entity.PermissionSet[action]
}

func (s *rbacService) EnforceForTargetAccount(ctx context.Context, tokenContext service_contract.TokenContext, userID string, targetAccountID string, action string) (bool, error) {
	if !s.isValidAction(action) {
		err := fmt.Errorf("invalid action permission format or name: %s", action)
		s.logger.Error(logger.Service, logger.RBACService, "failed to enforce policy", map[logger.ExtraKey]interface{}{
			logger.UserID:          userID,
			logger.TargetAccountID: targetAccountID,
			logger.Action:          action,
			logger.ErrorMessage:    err.Error(),
		})
		return false, err
	}

	allowed, err := s.repo.EnforceForTargetAccount(userID, targetAccountID, action)
	if err != nil {
		s.logger.Error(logger.Service, logger.RBACService, "failed to enforce policy", map[logger.ExtraKey]interface{}{
			logger.UserID:          userID,
			logger.TargetAccountID: targetAccountID,
			logger.Action:          action,
			logger.ErrorMessage:    err.Error(),
		})
		return false, err
	}
	return allowed, nil
}

func (s *rbacService) AddUserRoleForTargetAccount(ctx context.Context, tokenContext service_contract.TokenContext, userID string, targetAccountID string, role string) (bool, error) {
	if !s.isValidRole(role) {
		err := fmt.Errorf("invalid role: %s", role)
		s.logger.Error(logger.Service, logger.RBACService, "failed to add user role for target account", map[logger.ExtraKey]interface{}{
			logger.UserID:          userID,
			logger.Role:            role,
			logger.TargetAccountID: targetAccountID,
			logger.ErrorMessage:    err.Error(),
		})
		return false, err
	}

	added, err := s.repo.AddUserRoleForTargetAccount(userID, targetAccountID, role)
	if err != nil {
		s.logger.Error(logger.Service, logger.RBACService, "failed to add user role for target account", map[logger.ExtraKey]interface{}{
			logger.UserID:          userID,
			logger.Role:            role,
			logger.TargetAccountID: targetAccountID,
			logger.ErrorMessage:    err.Error(),
		})
		return false, err
	}
	s.logger.Info(logger.Service, logger.RBACService, "user role added for target account", map[logger.ExtraKey]interface{}{
		logger.UserID:          userID,
		logger.Role:            role,
		logger.TargetAccountID: targetAccountID,
	})
	return added, nil
}

func (s *rbacService) RemoveUserRoleForTargetAccount(ctx context.Context, tokenContext service_contract.TokenContext, userID string, targetAccountID string, role string) (bool, error) {
	if !s.isValidRole(role) {
		err := fmt.Errorf("invalid role: %s", role)
		s.logger.Error(logger.Service, logger.RBACService, "failed to remove user role for target account", map[logger.ExtraKey]interface{}{
			logger.UserID:          userID,
			logger.Role:            role,
			logger.TargetAccountID: targetAccountID,
			logger.ErrorMessage:    err.Error(),
		})
		return false, err
	}

	removed, err := s.repo.RemoveUserRoleForTargetAccount(userID, targetAccountID, role)
	if err != nil {
		s.logger.Error(logger.Service, logger.RBACService, "failed to remove user role for target account", map[logger.ExtraKey]interface{}{
			logger.UserID:          userID,
			logger.Role:            role,
			logger.TargetAccountID: targetAccountID,
			logger.ErrorMessage:    err.Error(),
		})
		return false, err
	}
	s.logger.Info(logger.Service, logger.RBACService, "user role removed for target account", map[logger.ExtraKey]interface{}{
		logger.UserID:          userID,
		logger.Role:            role,
		logger.TargetAccountID: targetAccountID,
	})
	return removed, nil
}

func (s *rbacService) GetUserRolesForTargetAccount(ctx context.Context, tokenContext service_contract.TokenContext, userID string, targetAccountID string) ([]string, error) {
	roles, err := s.repo.GetUserRolesForTargetAccount(userID, targetAccountID)
	if err != nil {
		s.logger.Error(logger.Service, logger.RBACService, "failed to get user roles for target account", map[logger.ExtraKey]interface{}{
			logger.UserID:          userID,
			logger.TargetAccountID: targetAccountID,
			logger.ErrorMessage:    err.Error(),
		})
		return nil, err
	}
	return roles, nil
}

func (s *rbacService) GetUsersForTargetAccount(ctx context.Context, tokenContext service_contract.TokenContext, targetAccountID string) ([][]string, error) {
	users, err := s.repo.GetUsersForTargetAccount(targetAccountID)
	if err != nil {
		s.logger.Error(logger.Service, logger.RBACService, "failed to get users for target account", map[logger.ExtraKey]interface{}{
			logger.TargetAccountID: targetAccountID,
			logger.ErrorMessage:    err.Error(),
		})
		return nil, err
	}
	return users, nil
}

func (s *rbacService) RemoveAllRolesForTargetAccount(ctx context.Context, tokenContext service_contract.TokenContext, targetAccountID string) (bool, error) {
	removed, err := s.repo.RemoveAllRolesForTargetAccount(targetAccountID)
	if err != nil {
		s.logger.Error(logger.Service, logger.RBACService, "failed to remove all roles for target account", map[logger.ExtraKey]interface{}{
			logger.TargetAccountID: targetAccountID,
			logger.ErrorMessage:    err.Error(),
		})
		return false, err
	}
	return removed, nil
}

func (s *rbacService) RemoveAllRolesForUser(ctx context.Context, tokenContext service_contract.TokenContext, userID string) (bool, error) {
	removed, err := s.repo.RemoveAllRolesForUser(userID)
	if err != nil {
		s.logger.Error(logger.Service, logger.RBACService, "failed to remove all roles for user", map[logger.ExtraKey]interface{}{
			logger.UserID:       userID,
			logger.ErrorMessage: err.Error(),
		})
		return false, err
	}
	return removed, nil
}

func (s *rbacService) AddPermissionForRole(ctx context.Context, tokenContext service_contract.TokenContext, role string, action string) (bool, error) {
	if !s.isValidRole(role) {
		err := fmt.Errorf("invalid role: %s", role)
		s.logger.Error(logger.Service, logger.RBACService, "failed to add permission for role", map[logger.ExtraKey]interface{}{
			logger.Role:         role,
			logger.Action:       action,
			logger.ErrorMessage: err.Error(),
		})
		return false, err
	}

	if !s.isValidAction(action) {
		err := fmt.Errorf("invalid action permission format or name: %s", action)
		s.logger.Error(logger.Service, logger.RBACService, "failed to add permission for role", map[logger.ExtraKey]interface{}{
			logger.Role:         role,
			logger.Action:       action,
			logger.ErrorMessage: err.Error(),
		})
		return false, err
	}

	added, err := s.repo.AddPermissionForRole(role, action)
	if err != nil {
		s.logger.Error(logger.Service, logger.RBACService, "failed to add permission for role", map[logger.ExtraKey]interface{}{
			logger.Role:         role,
			logger.Action:       action,
			logger.ErrorMessage: err.Error(),
		})
		return false, err
	}
	s.logger.Info(logger.Service, logger.RBACService, "permission added for role", map[logger.ExtraKey]interface{}{
		logger.Role:   role,
		logger.Action: action,
	})
	return added, nil
}

func (s *rbacService) RemovePermissionForRole(ctx context.Context, tokenContext service_contract.TokenContext, role string, action string) (bool, error) {
	if !s.isValidRole(role) {
		err := fmt.Errorf("invalid role: %s", role)
		s.logger.Error(logger.Service, logger.RBACService, "failed to remove permission for role", map[logger.ExtraKey]interface{}{
			logger.Role:         role,
			logger.Action:       action,
			logger.ErrorMessage: err.Error(),
		})
		return false, err
	}

	if !s.isValidAction(action) {
		err := fmt.Errorf("invalid action permission format or name: %s", action)
		s.logger.Error(logger.Service, logger.RBACService, "failed to remove permission for role", map[logger.ExtraKey]interface{}{
			logger.Role:         role,
			logger.Action:       action,
			logger.ErrorMessage: err.Error(),
		})
		return false, err
	}

	removed, err := s.repo.RemovePermissionForRole(role, action)
	if err != nil {
		s.logger.Error(logger.Service, logger.RBACService, "failed to remove permission for role", map[logger.ExtraKey]interface{}{
			logger.Role:         role,
			logger.Action:       action,
			logger.ErrorMessage: err.Error(),
		})
		return false, err
	}
	s.logger.Info(logger.Service, logger.RBACService, "permission removed for role", map[logger.ExtraKey]interface{}{
		logger.Role:   role,
		logger.Action: action,
	})
	return removed, nil
}

func (s *rbacService) GetPermissionsForRole(ctx context.Context, tokenContext service_contract.TokenContext, role string) ([][]string, error) {
	if !s.isValidRole(role) {
		err := fmt.Errorf("invalid role: %s", role)
		s.logger.Error(logger.Service, logger.RBACService, "failed to get permissions for role", map[logger.ExtraKey]interface{}{
			logger.Role:         role,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	perms, err := s.repo.GetPermissionsForRole(role)
	if err != nil {
		s.logger.Error(logger.Service, logger.RBACService, "failed to get permissions for role", map[logger.ExtraKey]interface{}{
			logger.Role:         role,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}
	return perms, nil
}
