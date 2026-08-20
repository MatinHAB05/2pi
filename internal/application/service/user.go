package service

import (
	"context"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/internal/domain/entity"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/pkg/logger"
)

type userService struct {
	repo                 repository_contract.UserRepository
	accountRepo          repository_contract.TargetAccountRepository
	rbacService          service_contract.RBACService
	userInfoCacheService service_contract.UserInfoCacheService
	logger               logger.Logger
}

func NewUserService(
	repo repository_contract.UserRepository,
	accountRepo repository_contract.TargetAccountRepository,
	rbacService service_contract.RBACService,
	log logger.Logger,
	userInfoCacheService service_contract.UserInfoCacheService,
) service_contract.UserService {
	return &userService{
		repo:                 repo,
		rbacService:          rbacService,
		accountRepo:          accountRepo,
		logger:               log,
		userInfoCacheService: userInfoCacheService,
	}
}

func (s *userService) Create(ctx context.Context, tokenContext service_contract.TokenContext, req service_contract.CreateUserRequest) (*service_contract.UserResponse, error) {
	user := &entity.User{
		BaseEntity: entity.BaseEntity{ID: int64(req.ID)},
		UserLang:   entity.LangFa,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Username:   req.Username,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to create user", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	s.logger.Info(logger.Service, logger.UserService, "user created successfully", map[logger.ExtraKey]interface{}{
		logger.UserID: user.ID,
	})

	return service_contract.ToUserResponse(user), nil
}

func (s *userService) GetByID(ctx context.Context, tokenContext service_contract.TokenContext, id int64) (*service_contract.UserResponse, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to get user by id", map[logger.ExtraKey]interface{}{
			logger.UserID:       id,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	return service_contract.ToUserResponse(u), nil
}

func (s *userService) GetWithTargetAccounts(ctx context.Context, tokenContext service_contract.TokenContext, id int64) (*service_contract.UserResponse, error) {
	u, err := s.repo.GetWithTargetAccounts(ctx, id)
	if err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to get user with target accounts", map[logger.ExtraKey]interface{}{
			logger.UserID:       id,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	return service_contract.ToUserResponse(u), nil
}

func (s *userService) GetByIDs(ctx context.Context, tokenContext service_contract.TokenContext, ids []int64) ([]*service_contract.UserResponse, error) {
	us, err := s.repo.GetByIDs(ctx, ids)
	if err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to get user by id", map[logger.ExtraKey]interface{}{
			logger.UserID + "s": ids,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	return service_contract.ToUserSliceResponse(us), nil
}

func (s *userService) GetByIDsWithTargetAccounts(ctx context.Context, tokenContext service_contract.TokenContext, ids []int64) ([]*service_contract.UserResponse, error) {
	us, err := s.repo.GetByIDsWithTargetAccounts(ctx, ids)
	if err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to get user with target accounts", map[logger.ExtraKey]interface{}{
			logger.UserID + "s": ids,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	return service_contract.ToUserSliceResponse(us), nil
}

func (s *userService) Update(ctx context.Context, tokenContext service_contract.TokenContext, req service_contract.UpdateUserRequest) (*service_contract.UserResponse, error) {
	user := &entity.User{
		BaseEntity: entity.BaseEntity{
			ID: int64(req.ID),
		},
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		UserLang:  req.Lang,
	}

	if err := s.repo.Update(ctx, user); err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to update user", map[logger.ExtraKey]interface{}{
			logger.UserID:       req.ID,
			"new_user_id":       user.ID,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	if err := s.userInfoCacheService.InvalidateCache(ctx, tokenContext, user.ID); err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to InvalidateCache user info", map[logger.ExtraKey]interface{}{
			logger.UserID:       req.ID,
			"new_user_id":       user.ID,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	s.logger.Info(logger.Service, logger.UserService, "user updated successfully", map[logger.ExtraKey]interface{}{
		logger.UserID: req.ID,
	})

	return service_contract.ToUserResponse(user), nil
}

func (s *userService) Delete(ctx context.Context, tokenContext service_contract.TokenContext, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to delete user", map[logger.ExtraKey]interface{}{
			logger.UserID:       id,
			logger.ErrorMessage: err.Error(),
		})
		return err
	}

	s.logger.Info(logger.Service, logger.UserService, "user deleted successfully", map[logger.ExtraKey]interface{}{
		logger.UserID: id,
	})
	return nil
}

func (s *userService) Exists(ctx context.Context, tokenContext service_contract.TokenContext, id int64) (bool, error) {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to check user existence", map[logger.ExtraKey]interface{}{
			logger.UserID:       id,
			logger.ErrorMessage: err.Error(),
		})
		return false, err
	}

	return exists, nil
}

func (s *userService) GetUserAccountsRolesByID(ctx context.Context, tokenContext service_contract.TokenContext, id int64) (*service_contract.UserAccountsRoleResponse, error) {
	u, err := s.rbacService.GetTargetAccountsForUser(ctx, service_contract.MapTokenContextToService(nil), id)
	if err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to get user-accounts-roles with user-id", map[logger.ExtraKey]interface{}{
			logger.UserID:       id,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	res, err := s.toUserAccountsRoleSliceResponse(ctx, u)
	if err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to build user accounts roles response", map[logger.ExtraKey]interface{}{
			logger.UserID:       id,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res[0], nil
}

func (s *userService) GetUserAccountsRolesByIDs(ctx context.Context, tokenContext service_contract.TokenContext, ids []int64) ([]*service_contract.UserAccountsRoleResponse, error) {
	u, err := s.rbacService.GetTargetAccountsForUsers(ctx, service_contract.MapTokenContextToService(nil), ids)
	if err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to get user-accounts-role with user-ids", map[logger.ExtraKey]interface{}{
			logger.UserID + "s": ids,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}
	s.logger.Info("", "", "", map[logger.ExtraKey]interface{}{
		"u": u,
	})
	res, err := s.toUserAccountsRoleSliceResponse(ctx, u)
	if err != nil {
		s.logger.Error(logger.Service, logger.UserService, "failed to build user accounts roles response", map[logger.ExtraKey]interface{}{
			logger.UserID + "s": ids,
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	return res, nil
}

type UserAccountRoleModelResponse struct {
	UserID          int64  `json:"user_id"`
	TargetAccountID int64  `json:"target_account_id"`
	Role            string `json:"role"`
}
type UserResponse struct {
	ID             int64       `json:"id"`
	Lang           entity.Lang `json:"lang"`
	FirstName      string
	LastName       string
	Username       string
	TargetAccounts []*TargetAccountResponse `json:"target_accounts,omitempty"`
}
type UserAccountsRoleResponse struct {
	*UserResponse
	AccountsRoles []*AccountsRoleResponse `json:",omitempty"`
}

type AccountsRoleResponse struct {
	AccountID *TargetAccountResponse
	Roles     []string
}

type TargetAccountResponse struct {
	ID          int64 `json:"id"`
	OwnerUserID int64 `json:"owner_user_id"`
	DayDuration int
	Period      int
	Description string
	Enable      bool
}

// toUserAccountsRoleResponse maps raw role entries for a single user into a structured UserAccountsRoleResponse.
func (s *userService) toUserAccountsRoleResponse(
	u []*service_contract.UserAccountRoleModelResponse,
	userMap map[int64]*service_contract.UserResponse,
	accMap map[int64]*service_contract.TargetAccountResponse,
) *service_contract.UserAccountsRoleResponse {
	if len(u) == 0 {
		return nil
	}

	userID := u[0].UserID
	userResp, exists := userMap[userID]
	if !exists {
		// Fallback if user details are missing from map
		userResp = &service_contract.UserResponse{ID: userID}
	}

	res := &service_contract.UserAccountsRoleResponse{
		UserResponse:  userResp,
		AccountsRoles: make([]*service_contract.AccountsRoleResponse, 0),
	}

	// Group roles by TargetAccountID for this specific user
	accountRolesMap := make(map[int64]*service_contract.AccountsRoleResponse)
	var accountOrder []int64

	for _, item := range u {
		if item == nil {
			continue
		}

		accRole, exists := accountRolesMap[item.TargetAccountID]
		if !exists {
			accResp, ok := accMap[item.TargetAccountID]
			if !ok {
				// Fallback if account details are missing from map
				accResp = &service_contract.TargetAccountResponse{ID: item.TargetAccountID}
			}

			accRole = &service_contract.AccountsRoleResponse{
				AccountID: accResp,
				Roles:     make([]string, 0),
			}
			accountRolesMap[item.TargetAccountID] = accRole
			accountOrder = append(accountOrder, item.TargetAccountID)
		}

		accRole.Roles = append(accRole.Roles, item.Role)
	}

	for _, accID := range accountOrder {
		res.AccountsRoles = append(res.AccountsRoles, accountRolesMap[accID])
	}

	return res
}

// toUserAccountsRoleSliceResponse fetches full user and account details, then groups roles per user/account.
func (s *userService) toUserAccountsRoleSliceResponse(
	ctx context.Context,
	u []*service_contract.UserAccountRoleModelResponse,
) ([]*service_contract.UserAccountsRoleResponse, error) {
	if len(u) == 0 {
		return nil, nil
	}

	// Step 1: Collect unique user IDs and account IDs
	userIDsMap := make(map[int64]struct{})
	accIDsMap := make(map[int64]struct{})

	for _, item := range u {
		if item == nil {
			continue
		}
		userIDsMap[item.UserID] = struct{}{}
		accIDsMap[item.TargetAccountID] = struct{}{}
	}

	userIDs := make([]int64, 0, len(userIDsMap))
	for id := range userIDsMap {
		userIDs = append(userIDs, id)
	}

	accIDs := make([]int64, 0, len(accIDsMap))
	for id := range accIDsMap {
		accIDs = append(accIDs, id)
	}

	// Step 2: Batch fetch full user and account entities
	users, err := s.repo.GetByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	accounts, err := s.accountRepo.GetByIDs(ctx, accIDs)
	if err != nil {
		return nil, err
	}

	// Step 3: Build lookup maps for fast O(1) access
	userMap := make(map[int64]*service_contract.UserResponse, len(users))
	for i := range users {
		userMap[users[i].ID] = service_contract.ToUserResponse(&users[i])
	}

	accMap := make(map[int64]*service_contract.TargetAccountResponse, len(accounts))
	for i := range accounts {
		accMap[accounts[i].ID] = service_contract.ToTargetAccountResponse(&accounts[i])
	}

	// Step 4: Group raw entries by UserID
	groupedByUser := make(map[int64][]*service_contract.UserAccountRoleModelResponse)
	var userOrder []int64

	for _, item := range u {
		if item == nil {
			continue
		}
		if _, exists := groupedByUser[item.UserID]; !exists {
			userOrder = append(userOrder, item.UserID)
		}
		groupedByUser[item.UserID] = append(groupedByUser[item.UserID], item)
	}

	// Step 5: Construct final response slice with populated user and account details
	final := make([]*service_contract.UserAccountsRoleResponse, 0, len(userOrder))
	for _, userID := range userOrder {
		items := groupedByUser[userID]
		if resp := s.toUserAccountsRoleResponse(items, userMap, accMap); resp != nil {
			final = append(final, resp)
		}
	}

	return final, nil
}
