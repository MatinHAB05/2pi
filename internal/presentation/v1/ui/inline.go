package ui

import (
	"fmt"
	"strconv"

	"github.com/go-telegram/bot/models"
)

// OnboardingInlineKeyboard guides new users to complete their profile setup
func OnboardingInlineKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "🚀 Complete Account Setup", CallbackData: "acc:edit_menu"},
			},
		},
	}
}

// ActiveAccountInlineKeyboard renders context actions based on RBAC role
func ActiveAccountInlineKeyboard(role string, enabled bool) *models.InlineKeyboardMarkup {
	statusText := "⚡ Enable Account"
	if enabled {
		statusText = "⚡ Disable Account"
	}

	if role == "viewer" {
		return &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "🔄 Switch Account", CallbackData: "switch:menu"},
				},
			},
		}
	}

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: statusText, CallbackData: "acc:toggle_status"},
				{Text: "✏️ Edit Details", CallbackData: "acc:edit_menu"},
			},
			{
				{Text: "🔄 Switch Account", CallbackData: "switch:menu"},
			},
		},
	}
}

// SwitchAccountInlineKeyboard dynamically lists target accounts user has access to
func SwitchAccountInlineKeyboard(accounts []AccountItem, activeID int64) *models.InlineKeyboardMarkup {
	var rows [][]models.InlineKeyboardButton

	for _, acc := range accounts {
		prefix := "⚪"
		suffix := ""
		if acc.ID == activeID {
			prefix = "🟢"
			suffix = " [ACTIVE]"
		}
		label := fmt.Sprintf("%s @%s (%s)%s", prefix, acc.Username, acc.Role, suffix)
		rows = append(rows, []models.InlineKeyboardButton{
			{Text: label, CallbackData: "switch:select:" + strconv.FormatInt(acc.ID, 10)},
		})
	}

	rows = append(rows,
		[]models.InlineKeyboardButton{{Text: "➕ Create New Target Account", CallbackData: "acc:create"}},
		[]models.InlineKeyboardButton{{Text: "🚪 Logout (Clear Selection)", CallbackData: "switch:logout"}},
	)

	return &models.InlineKeyboardMarkup{InlineKeyboard: rows}
}

// EditAccountInlineKeyboard renders target account parameters
func EditAccountInlineKeyboard(enabled bool) *models.InlineKeyboardMarkup {
	statusLabel := "⚡ Status: [ Disabled 🔴 ]"
	if enabled {
		statusLabel = "⚡ Status: [ Enabled 🟢 ]"
	}

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "✏️ Username", CallbackData: "edit:field:username"},
				{Text: "⏱ Day Duration", CallbackData: "edit:field:day_duration"},
			},
			{
				{Text: "🔄 Period", CallbackData: "edit:field:period"},
				{Text: "🌐 User Language", CallbackData: "edit:field:user_lang"},
			},
			{
				{Text: "📝 Description", CallbackData: "edit:field:description"},
			},
			{
				{Text: statusLabel, CallbackData: "acc:toggle_status"},
			},
			{
				{Text: "🔙 Back to Dashboard", CallbackData: "acc:dashboard"},
			},
		},
	}
}

// ShareAccessInlineKeyboard options for access management
func ShareAccessInlineKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "➕ Invite New User", CallbackData: "share:invite"},
			},
			{
				{Text: "📋 List Active Collaborators", CallbackData: "share:list"},
			},
			{
				{Text: "🔙 Back to Dashboard", CallbackData: "acc:dashboard"},
			},
		},
	}
}

// SelectRoleInlineKeyboard selects RBAC permissions for invitations
func SelectRoleInlineKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "🛠 Admin", CallbackData: "invite:role:admin"},
				{Text: "✏️ Editor", CallbackData: "invite:role:editor"},
			},
			{
				{Text: "👁 Viewer", CallbackData: "invite:role:viewer"},
			},
			{
				{Text: "❌ Cancel", CallbackData: "share:menu"},
			},
		},
	}
}

// AcceptInviteInlineKeyboard renders invitation acceptance prompt
func AcceptInviteInlineKeyboard(token string) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "✅ Accept Invitation", CallbackData: "invite:accept:" + token},
				{Text: "❌ Decline", CallbackData: "invite:decline:" + token},
			},
		},
	}
}
