package handler

import (
	"fmt"
	"log"
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
	MsgSettingsFormat      = "Bot settings for user ID: %d\nLanguage : %s"
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

func WhoCanAccessToThisAccountMsg(users []*service_contract.UserAccountsRoleResponse) string {
	if len(users) == 0 {
		return "No users have access to this account."
	}

	var sb strings.Builder
	sb.WriteString("Users with access to this account:\n")

	for i := range users {
		if users[i] == nil {
			continue
		}
		sb.WriteString(ShowUserInfo(users[i].UserResponse))
		if len(users[i].AccountsRoles) != 1 {
			log.Println("WTF!!!! - check `WhoCanAccessToThisAccountMsg`")
		}
		for acc := range users[i].AccountsRoles { // MUST BE 1 ACCOUNT!!!
			for _, role := range users[i].AccountsRoles[acc].Roles {
				sb.WriteString(fmt.Sprintf("• Role : `%s`", role))
			}
		}
		sb.WriteString("\n\n---------------\n\n")

	}

	return sb.String()
}

func ShowCurrentAccountInfo(user *service_contract.UserResponse, acc *service_contract.TargetAccountResponse) string {
	if user == nil || acc == nil {
		return "User or account details not found."
	}

	var sb strings.Builder

	// User Information Section
	sb.WriteString(ShowUserInfo(user))

	sb.WriteString("\n\n-------------------------\n\n")

	// Active Account Section
	sb.WriteString(ShowAccountInfo(acc))

	return sb.String()
}

func ShowUserInfo(user *service_contract.UserResponse) string {
	if user == nil {
		return "User or account details not found."
	}

	var sb strings.Builder

	// User Information Section
	sb.WriteString("👤 **User Info**\n")
	sb.WriteString(fmt.Sprintf("• User ID: `%d`\n", user.ID))
	sb.WriteString(fmt.Sprintf("• User FirstName: `%s`\n", user.FirstName))
	sb.WriteString(fmt.Sprintf("• User LastName: `%s`\n", user.LastName))
	sb.WriteString(fmt.Sprintf("• User Username: `%s`\n", user.Username))
	sb.WriteString(fmt.Sprintf("• Language: `%s`\n", user.Lang))
	return sb.String()
}

func ShowAccountInfo(acc *service_contract.TargetAccountResponse) string {
	if acc == nil {
		return "User or account details not found."
	}

	var sb strings.Builder

	// User Information Section
	sb.WriteString("🔑 **Active Account Info**\n")
	sb.WriteString(fmt.Sprintf("• Account ID: `%d`\n", acc.ID))
	sb.WriteString(fmt.Sprintf("• Owner ID: `%d`\n", acc.OwnerUserID))
	sb.WriteString(fmt.Sprintf("• Duration: `%d` Days\n", acc.DayDuration))
	sb.WriteString(fmt.Sprintf("• Period: `%d`\n", acc.Period))
	sb.WriteString(fmt.Sprintf("• Status: `%t`\n", acc.Enable))
	if acc.Description != "" {
		sb.WriteString(fmt.Sprintf("• Description: %s\n", acc.Description))
	}

	return sb.String()
}
