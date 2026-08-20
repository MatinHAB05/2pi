package handler

import (
	"sync"

	"github.com/MatinHAB05/2pi/internal/presentation/v1/ui"
)

// ============================================================================
// IN-MEMORY MOCK DATABASE
// ============================================================================

type MockDB struct {
	mu             sync.RWMutex
	knownUsers     map[int64]bool
	currentContext map[int64]int64
	userAccounts   map[int64]map[int64]*ui.AccountItem
	invitations    map[string]Invitation
}

var db = &MockDB{
	knownUsers:     make(map[int64]bool),
	currentContext: make(map[int64]int64),
	userAccounts:   make(map[int64]map[int64]*ui.AccountItem),
	invitations:    make(map[string]Invitation),
}

// func initMockAccountsForUser(chatID int64) {
// 	db.userAccounts[chatID] = map[int64]*ui.AccountItem{
// 		101: {
// 			ID:          101,
// 			Username:    "main_entry_tracker",
// 			Role:        "owner",
// 			DayDuration: 7,
// 			Period:      30,
// 			UserLang:    "fa",
// 			Description: "Primary personal entry period cycle",
// 			Enable:      true,
// 		},
// 		102: {
// 			ID:          102,
// 			Username:    "dev_server_renewals",
// 			Role:        "editor",
// 			DayDuration: 14,
// 			Period:      90,
// 			UserLang:    "eng",
// 			Description: "Team infrastructure entry check",
// 			Enable:      true,
// 		},
// 	}
// 	db.currentContext[chatID] = 101
// }

// func checkUserIsNew(id int64) bool {
// 	db.mu.RLock()
// 	defer db.mu.RUnlock()
// 	return !db.knownUsers[id]
// }

// func createDefaultAccount(id int64) {
// 	db.mu.Lock()
// 	defer db.mu.Unlock()

// 	db.knownUsers[id] = true
// 	db.userAccounts[id] = map[int64]*ui.AccountItem{
// 		201: {
// 			ID:          201,
// 			Username:    "new_user_account",
// 			Role:        "owner",
// 			DayDuration: 0,
// 			Period:      0,
// 			UserLang:    "fa",
// 			Description: "Unconfigured default account",
// 			Enable:      false,
// 		},
// 	}
// 	db.currentContext[id] = 201
// }

// func getCurrentAccountID(id int64) int64 {
// 	db.mu.RLock()
// 	defer db.mu.RUnlock()
// 	return db.currentContext[id]
// }

// func getCurrentAccount(id int64) ui.AccountItem {
// 	db.mu.Lock()
// 	if _, exists := db.userAccounts[id]; !exists {
// 		if db.knownUsers[id] {
// 			initMockAccountsForUser(id)
// 		} else {
// 			db.mu.Unlock()
// 			createDefaultAccount(id)
// 			db.mu.Lock()
// 		}
// 	}
// 	db.mu.Unlock()

// 	db.mu.RLock()
// 	defer db.mu.RUnlock()

// 	currentID := db.currentContext[id]
// 	if acc, exists := db.userAccounts[id][currentID]; exists {
// 		return *acc
// 	}

// 	return ui.AccountItem{
// 		Username:    "unknown_account",
// 		Role:        "owner",
// 		UserLang:    "fa",
// 		Description: "No current context found",
// 	}
// }

// func getUserAccessibleAccounts(id int64) []ui.AccountItem {
// 	db.mu.RLock()
// 	defer db.mu.RUnlock()

// 	var list []ui.AccountItem
// 	for _, acc := range db.userAccounts[id] {
// 		list = append(list, *acc)
// 	}
// 	return list
// }

// func setCurrentAccountContext(chatID int64, accIDStr string) {
// 	accID, err := strconv.ParseInt(accIDStr, 10, 64)
// 	if err != nil {
// 		return
// 	}
// 	db.mu.Lock()
// 	defer db.mu.Unlock()
// 	db.currentContext[chatID] = accID
// }

// func clearCurrentAccountContext(chatID int64) {
// 	db.mu.Lock()
// 	defer db.mu.Unlock()
// 	delete(db.currentContext, chatID)
// }

// func generateInviteLink(chatID int64, role string) string {
// 	db.mu.Lock()
// 	defer db.mu.Unlock()

// 	currentAcc := getCurrentAccount(chatID)
// 	token := fmt.Sprintf("tok_%d_%s", currentAcc.ID, role)
// 	db.invitations[token] = Invitation{
// 		TargetUsername: currentAcc.Username,
// 		Role:           role,
// 	}
// 	return fmt.Sprintf("https://t.me/YourPeriodBot?start=invite_%s", token)
// }

// func getInvitationByToken(token string) Invitation {
// 	db.mu.RLock()
// 	defer db.mu.RUnlock()

// 	if inv, exists := db.invitations[token]; exists {
// 		return inv
// 	}
// 	return Invitation{TargetUsername: "shared_team_account", Role: "editor"}
// }

// ============================================================================
// HELPER UTILITIES
// ============================================================================

// func getChatID(u *models.Update) int64 {
// 	if u.Message != nil {
// 		return u.Message.Chat.ID
// 	}
// 	if u.CallbackQuery != nil && u.CallbackQuery.Message.Message != nil {
// 		return u.CallbackQuery.Message.Message.Chat.ID
// 	}
// 	return 0
// }

// func sanitizeMarkdownText(acc ui.AccountItem) (string, string, string, string) {
// 	username := acc.Username
// 	if username == "" {
// 		username = "<unnamed>"
// 	}
// 	role := acc.Role
// 	if role == "" {
// 		role = "owner"
// 	}
// 	lang := string(acc.UserLang)
// 	if lang == "" {
// 		lang = "fa"
// 	}
// 	desc := acc.Description
// 	if desc == "" {
// 		desc = "No description configured"
// 	}
// 	return username, role, lang, desc
// }

// func formatStatus(enable bool) string {
// 	if enable {
// 		return "🟢 Current"
// 	}
// 	return "🔴 Disabled"
// }

type Invitation struct {
	TargetUsername string
	Role           string
}
