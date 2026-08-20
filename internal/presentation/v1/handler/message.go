package handler

import (
	"fmt"
	"strings"

	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
)

const (
	// Callback Prefixes & States
	AccountFieldHandlerPrefix = "acc:edit:fields:field:handler:"
	StateEditFieldsEnter      = "acc:edit:fields:field:enter"

	UserChangeLanguageHandlerPrefix = "user:fields:field:language:"

	ShareAccountHandlerPrefix               = "acc:share:handler:"
	InviteAccountHandlerPrefix              = "acc:share:invite:handler:role:"
	StateConfirmShareAccountAccessCodeEnter = "acc:confirm:share:enter"

	SwitchCurrentAccountHandlerPrefix = "acc:switch:select:handler:"

	// State Keys
	StateKeyStatus = "state"
	StateKeyAccID  = "acc_id"
	StateKeyField  = "field"

	// Field Identifiers
	AccountFieldDayDuration  = "day_duration"
	AccountFieldPeriod       = "period"
	AccountFieldDescription  = "description"
	AccountFieldToggleStatus = "toggle_status"
	AccountFieldDashboard    = "dashboard"

	ShareAccountAccessInvite        = "invite"
	ShareAccountAccessConfirmInvite = "confirm-invite"
	ShareAccountAccessList          = "list"
	ShareAccountAccessDashboard     = "dashboard"

	InviteAccountAccessRoleAdmin  = "admin"
	InviteAccountAccessRoleOwner  = "owner"
	InviteAccountAccessRoleEditor = "editor"
	InviteAccountAccessRoleViewer = "viewer"
	InviteAccountAccessRoleCancel = "cancel"

	UserChangeLanguageHandlerFarsi   = "fa"
	UserChangeLanguageHandlerEnglish = "eng"
	UserChangeLanguageHandlerCancel  = "cancel"

	// Response & Navigation Messages
	MsgWelcome             = "👋 Welcome to Period Tracker Bot!"
	MsgEditAccountSettings = "⚙️ Edit Account Settings \n\nSelect a parameter to update:"
	MsgHelp                = "Need help? Here are the available commands..."
	MsgSettingsFormat      = "Bot settings for user ID: %d"
	MsgUnknownCommand      = "unknown : /help"
	MsgSuccessDone         = "✅ Done"
	MsgSwitchAccount       = "Switch Account Message :"
	MsgShareAccountMenu    = "Share Account Message Menu :"
	MsgShareAccountInvite  = "Share Account Invite Message"
	MsgChangeLanguage      = "Language Changed :\nNew : %s\nReq : %s"
	MsgChangeLanguageMenu  = "Change Language"

	// Input Prompts
	MsgEnterNewDayDuration                = "Enter New duration in days :"
	MsgEnterNewPeriod                     = "Enter New period in days :"
	MsgEnterNewDescription                = "Enter New Description :"
	MsgEnterOneTimeShareAccountAccessCode = "Enter One-Time Share Account Access Code :"

	// Error Messages
	MsgErrInvalidDayDuration  = "❌ Invalid day duration value. Please enter a valid number."
	MsgErrInvalidPeriod       = "❌ Invalid period value. Please enter a valid number."
	MsgErrUpdateAccountFailed = "❌ Failed to update account settings. Please try again later."
)

func WhoCanAccessToThisAccountMsg(users []*service_contract.UserResponse) string {
	if len(users) == 0 {
		return "No users have access to this account."
	}

	var sb strings.Builder
	sb.WriteString("Users with access to this account:\n")

	for _, user := range users {
		if user == nil {
			continue
		}
		sb.WriteString(fmt.Sprintf("• User ID: %d\n", user.ID))
	}

	return sb.String()
}
