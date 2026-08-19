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
				{Text: "🚀 Account Setup", CallbackData: "acc:edit_menu"},
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

func CancelEditInlineKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "❌ Cancel", CallbackData: "edit:cancel"},
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
		label := fmt.Sprintf("%s (%s)%s", prefix, acc.Role, suffix)
		rows = append(rows, []models.InlineKeyboardButton{
			{Text: label, CallbackData: "switch:select:" + strconv.FormatInt(acc.ID, 10)},
		})
	}

	return &models.InlineKeyboardMarkup{InlineKeyboard: rows}
}

// EditAccountInlineKeyboard renders target account parameters
func EditAccountInlineKeyboard(enabled bool) *models.InlineKeyboardMarkup {
	statusLabel := "🟢 Enabled"
	if enabled {
		statusLabel = "🔴 Disabled"
	}

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "⏱ Day Duration", CallbackData: "acc:edit:fields:field:handler:day_duration"},
				{Text: "🔄 Period", CallbackData: "acc:edit:fields:field:handler:period"},
			},
			{
				{Text: statusLabel, CallbackData: "acc:edit:fields:field:handler:toggle_status"},
				{Text: "📝 Description", CallbackData: "acc:edit:fields:field:handler:description"},
			},
			{
				{Text: "🔙 Back to Dashboard", CallbackData: "acc:edit:fields:field:handler:dashboard"},
			},
		},
	}
}

// ShareAccessInlineKeyboard options for access management
func ShareAccessInlineKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "➕ Invite Access", CallbackData: "acc:share:handler:invite"},
				{Text: "🫂 Confirm Invite", CallbackData: "acc:share:handler:confirm-invite"},
			},
			{
				{Text: "📋 List Active Collaborators", CallbackData: "acc:share:handler:list"},
			},
			{
				{Text: "🔙 Back to Dashboard", CallbackData: "acc:share:handler:dashboard"},
			},
		},
	}
}

// SelectRoleInlineKeyboard selects RBAC permissions for invitations
func SelectRoleInlineKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "🛠 Admin", CallbackData: "acc:share:invite:handler:role:admin"},
				{Text: "🦉 Owner", CallbackData: "acc:share:invite:handler:role:owner"},
			},
			{
				{Text: "✏️ Editor", CallbackData: "acc:share:invite:handler:role:editor"},
				{Text: "👁 Viewer", CallbackData: "acc:share:invite:handler:role:viewer"},
			},
			{
				{Text: "❌ Cancel", CallbackData: "acc:share:invite:handler:role:cancel"},
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

// OnboardingInlineKeyboard guides new users to complete their profile setup
func ChangeLanguageInlineKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "🇮🇷 Farsi", CallbackData: "user:fields:field:language:fa"},
				{Text: "🇺🇸 English", CallbackData: "user:fields:field:language:eng"},
			},
			{
				{Text: "❌ Cancel", CallbackData: "user:fields:field:language:cancel"},
			},
		},
	}
}
