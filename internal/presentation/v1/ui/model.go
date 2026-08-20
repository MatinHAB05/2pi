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

func MapTargetAccountToAccountItem(account *service_contract.TargetAccountResponse) *AccountItem {
	var res *AccountItem
	if account == nil {
		return res
	}
	res = &AccountItem{
		ID:          account.ID,
		Role:        "<!Hi!>",
		DayDuration: account.DayDuration,
		Period:      account.Period,
		Description: account.Description,
		Enable:      account.Enable,
	}
	return res
}

func MapTargetAccountsToAccountItems(accounts []*service_contract.TargetAccountResponse) []*AccountItem {
	res := make([]*AccountItem, len(accounts))
	if len(accounts) == 0 {
		return res
	}
	for i := range accounts {
		res[i] = MapTargetAccountToAccountItem(accounts[i])
	}
	return res
}
