package handler

const (
	// Callback Prefixes & States
	AccountFieldHandlerPrefix  = "acc:edit:fields:field:handler:"
	StateEditFieldsEnter       = "acc:edit:fields:field:enter"
	ShareAccountHandlerPrefix  = "acc:share:handler:"
	InviteAccountHandlerPrefix = "acc:share:invite:handler:role:"

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

	ShareAccountAccessInvite    = "invite"
	ShareAccountAccessList      = "list"
	ShareAccountAccessDashboard = "dashboard"

	InviteAccountAccessRoleAdmin  = "admin"
	InviteAccountAccessRoleOwner  = "owner"
	InviteAccountAccessRoleEditor = "editor"
	InviteAccountAccessRoleViewer = "viewer"
	InviteAccountAccessRoleCancel = "cancel"

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

	// Input Prompts
	MsgEnterNewDayDuration = "Enter New duration in days :"
	MsgEnterNewPeriod      = "Enter New period in days :"
	MsgEnterNewDescription = "Enter New Description :"

	// Error Messages
	MsgErrInvalidDayDuration  = "❌ Invalid day duration value. Please enter a valid number."
	MsgErrInvalidPeriod       = "❌ Invalid period value. Please enter a valid number."
	MsgErrUpdateAccountFailed = "❌ Failed to update account settings. Please try again later."
)
