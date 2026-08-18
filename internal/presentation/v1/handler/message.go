package handler

const (
	// Callback Prefixes & States
	FieldHandlerPrefix   = "acc:edit:fields:field:handler:"
	StateEditFieldsEnter = "acc:edit:fields:field:enter"

	// State Keys
	StateKeyStatus = "state"
	StateKeyAccID  = "acc_id"
	StateKeyField  = "field"

	// Field Identifiers
	FieldDayDuration  = "day_duration"
	FieldPeriod       = "period"
	FieldDescription  = "description"
	FieldToggleStatus = "toggle_status"
	FieldDashboard    = "dashboard"

	// Response & Navigation Messages
	MsgWelcome             = "👋 Welcome to Period Tracker Bot!"
	MsgDashboardWelcome    = "👋 Welcome to Period Tracker Bot!\n\nWe created your default account profile. Tracking is disabled until setup is completed."
	MsgEditAccountSettings = "⚙️ Edit Account Settings \n\nSelect a parameter to update:"
	MsgHelp                = "Need help? Here are the available commands..."
	MsgSettingsFormat      = "Bot settings for user ID: %d"
	MsgUnknownCommand      = "unknown : /help"
	MsgSuccessDone         = "✅ Done"

	// Input Prompts
	MsgEnterNewDayDuration = "Enter New duration in days :"
	MsgEnterNewPeriod      = "Enter New period in days :"
	MsgEnterNewDescription = "Enter New Description :"

	// Error Messages
	MsgErrInvalidDayDuration  = "❌ Invalid day duration value. Please enter a valid number."
	MsgErrInvalidPeriod       = "❌ Invalid period value. Please enter a valid number."
	MsgErrUpdateAccountFailed = "❌ Failed to update account settings. Please try again later."
)
