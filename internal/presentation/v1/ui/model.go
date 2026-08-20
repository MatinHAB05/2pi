package ui

import service_contract "github.com/MatinHAB05/2pi/internal/application/contract"

type AccountItem struct {
	ID          int64
	Role        string
	DayDuration int
	Period      int
	Description string
	Enable      bool
}

func MapTargetAccountToAccountItem(account *service_contract.AccountsRoleResponse) []*AccountItem {
	if account == nil || account.Account == nil || len(account.Roles) == 0 {
		return []*AccountItem{}
	}

	res := make([]*AccountItem, len(account.Roles))
	for i, role := range account.Roles {
		res[i] = &AccountItem{
			ID:          account.Account.ID,
			Role:        role,
			DayDuration: account.Account.DayDuration,
			Period:      account.Account.Period,
			Description: account.Account.Description,
			Enable:      account.Account.Enable,
		}
	}
	return res
}

func MapTargetAccountsToAccountItems(accounts *service_contract.UserAccountsRoleResponse) []*AccountItem {
	if accounts == nil || len(accounts.AccountsRoles) == 0 {
		return []*AccountItem{}
	}

	res := make([]*AccountItem, 0, len(accounts.AccountsRoles))
	for _, accRole := range accounts.AccountsRoles {
		res = append(res, MapTargetAccountToAccountItem(accRole)...)
	}

	return res
}
