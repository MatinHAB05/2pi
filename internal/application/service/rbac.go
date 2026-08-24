package service

import (
	"context"
	"fmt"
	"strconv"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/entity"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/internal/infrastructure/database"
	"github.com/MatinHAB05/2pi/pkg/logger"
)

// TODO : seed = action-role - Done
type rbacService struct {
	repo       repository_contract.RBACRepository
	logger     logger.Logger
	trxManager database.TrxManager
}

func NewRBACService(
	repo repository_contract.RBACRepository,
	log logger.Logger,
	trxManager database.TrxManager,
) service_contract.RBACService {
	return &rbacService{
		repo:       repo,
		logger:     log,
		trxManager: trxManager,
	}
}
func (s *rbacService) isValidRole(role string) bool {
	return entity.RoleSet[role]
}

func (s *rbacService) isValidAction(action string) bool {
	return entity.PermissionSet[action]
}

func (s *rbacService) EnforceForTargetAccount(ctx context.Context, tokenContext service_contract.TokenContext, userID int64, targetAccountID int64, action string) (bool, error) {
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

	userIDStr := strconv.FormatInt(userID, 10)
	targetAccountIDStr := strconv.FormatInt(targetAccountID, 10)

	allowed, err := s.repo.EnforceForTargetAccount(ctx, userIDStr, targetAccountIDStr, action)
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

func (s *rbacService) AddUserRoleForTargetAccount(ctx context.Context, tokenContext service_contract.TokenContext, userID int64, targetAccountID int64, role string) (bool, error) {
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

	userIDStr := strconv.FormatInt(userID, 10)
	targetAccountIDStr := strconv.FormatInt(targetAccountID, 10)

	added, err := s.repo.AddUserRoleForTargetAccount(ctx, userIDStr, targetAccountIDStr, role)
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

func (s *rbacService) RemoveUserRoleForTargetAccount(ctx context.Context, tokenContext service_contract.TokenContext, userID int64, targetAccountID int64, role string) (bool, error) {
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

	userIDStr := strconv.FormatInt(userID, 10)
	targetAccountIDStr := strconv.FormatInt(targetAccountID, 10)

	removed, err := s.repo.RemoveUserRoleForTargetAccount(ctx, userIDStr, targetAccountIDStr, role)
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

func (s *rbacService) GetUserRolesForTargetAccount(ctx context.Context, tokenContext service_contract.TokenContext, userID int64, targetAccountID int64) ([]*service_contract.UserAccountRoleModelResponse, error) {
	userIDStr := strconv.FormatInt(userID, 10)
	targetAccountIDStr := strconv.FormatInt(targetAccountID, 10)

	roles, err := s.repo.GetUserRolesForTargetAccount(ctx, userIDStr, targetAccountIDStr)
	if err != nil {
		s.logger.Error(logger.Service, logger.RBACService, "failed to get user roles for target account", map[logger.ExtraKey]interface{}{
			logger.UserID:          userID,
			logger.TargetAccountID: targetAccountID,
			logger.ErrorMessage:    err.Error(),
		})
		return nil, err
	}
	return service_contract.ToUserAccountRoleSliceModelResponse(roles), nil
}

func (s *rbacService) GetUsersForTargetAccount(ctx context.Context, tokenContext service_contract.TokenContext, targetAccountID int64) ([]*service_contract.UserAccountRoleModelResponse, error) {
	targetAccountIDStr := strconv.FormatInt(targetAccountID, 10)

	users, err := s.repo.GetUsersForTargetAccount(ctx, targetAccountIDStr)
	if err != nil {
		s.logger.Error(logger.Service, logger.RBACService, "failed to get users for target account", map[logger.ExtraKey]interface{}{
			logger.TargetAccountID: targetAccountID,
			logger.ErrorMessage:    err.Error(),
		})
		return nil, err
	}
	return service_contract.ToUserAccountRoleSliceModelResponse(users), nil
}

func (s *rbacService) GetTargetAccountsForUser(ctx context.Context, tokenContext service_contract.TokenContext, userID int64) ([]*service_contract.UserAccountRoleModelResponse, error) {
	userIDStr := strconv.FormatInt(userID, 10)

	accs, err := s.repo.GetTargetAccountsForUser(ctx, userIDStr)
	if err != nil {
		s.logger.Error(logger.Service, logger.RBACService, "failed to get accounts for user", map[logger.ExtraKey]interface{}{
			logger.UserID:       userID,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}
	return service_contract.ToUserAccountRoleSliceModelResponse(accs), nil
}

func (s *rbacService) GetTargetAccountsForUsers(ctx context.Context, tokenContext service_contract.TokenContext, userIDs []int64) ([]*service_contract.UserAccountRoleModelResponse, error) {
	strIDs := make([]string, 0, len(userIDs))
	for _, i := range userIDs {
		strIDs = append(strIDs, strconv.FormatInt(i, 10))
	}
	fmt.Println(strIDs)
	accs, err := s.repo.GetTargetAccountsForUsers(ctx, strIDs)
	if err != nil {
		s.logger.Error(logger.Service, logger.RBACService, "failed to get accounts for users", map[logger.ExtraKey]interface{}{
			logger.UserID + "s": userIDs,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}
	s.logger.Info("", "", "", map[logger.ExtraKey]interface{}{
		"accs": accs,
	})
	s.logger.Info("", "", "", map[logger.ExtraKey]interface{}{
		"Map": service_contract.ToUserAccountRoleSliceModelResponse(accs),
	})
	return service_contract.ToUserAccountRoleSliceModelResponse(accs), nil
}

func (s *rbacService) RemoveAllRolesForTargetAccount(ctx context.Context, tokenContext service_contract.TokenContext, targetAccountID int64) (bool, error) {
	targetAccountIDStr := strconv.FormatInt(targetAccountID, 10)

	removed, err := s.repo.RemoveAllRolesForTargetAccount(ctx, targetAccountIDStr)
	if err != nil {
		s.logger.Error(logger.Service, logger.RBACService, "failed to remove all roles for target account", map[logger.ExtraKey]interface{}{
			logger.TargetAccountID: targetAccountID,
			logger.ErrorMessage:    err.Error(),
		})
		return false, err
	}
	return removed, nil
}

func (s *rbacService) RemoveAllRolesForUser(ctx context.Context, tokenContext service_contract.TokenContext, userID int64) (bool, error) {
	userIDStr := strconv.FormatInt(userID, 10)

	removed, err := s.repo.RemoveAllRolesForUser(ctx, userIDStr)
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

	added, err := s.repo.AddPermissionForRole(ctx, role, action)
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

	removed, err := s.repo.RemovePermissionForRole(ctx, role, action)
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

func (s *rbacService) GetPermissionsForRole(ctx context.Context, tokenContext service_contract.TokenContext, role string) ([]*service_contract.RolePermissionModelResponse, error) {
	if !s.isValidRole(role) {
		err := fmt.Errorf("invalid role: %s", role)
		s.logger.Error(logger.Service, logger.RBACService, "failed to get permissions for role", map[logger.ExtraKey]interface{}{
			logger.Role:         role,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	perms, err := s.repo.GetPermissionsForRole(ctx, role)
	if err != nil {
		s.logger.Error(logger.Service, logger.RBACService, "failed to get permissions for role", map[logger.ExtraKey]interface{}{
			logger.Role:         role,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}
	return service_contract.ToRolePermissionSliceModelModelResponse(perms), nil
}
