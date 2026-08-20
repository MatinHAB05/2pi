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
	res := make([]*AccountItem, len(account.Roles))
	if len(account.Roles) == 0 {
		return res
	}
	for i := range account.Roles {
		res[i] = &AccountItem{
			ID:          account.AccountID.ID,
			Role:        account.Roles[i],
			DayDuration: account.AccountID.DayDuration,
			Period:      account.AccountID.Period,
			Description: account.AccountID.Description,
			Enable:      account.AccountID.Enable,
		}
	}
	return res
}

func MapTargetAccountsToAccountItems(accounts *service_contract.UserAccountsRoleResponse) []*AccountItem {
	res := make([]*AccountItem, len(accounts.AccountsRoles)) // at least!
	if len(accounts.AccountsRoles) == 0 {
		return res
	}
	for i := range accounts.AccountsRoles {
		res = append(res, MapTargetAccountToAccountItem(accounts.AccountsRoles[i])...)
	}

	return res
}
