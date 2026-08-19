package entity

type Role string
type Resource string
type Action string

// Roles Definitions
const (
	RoleAdmin  Role = "admin"
	RoleOwner  Role = "owner"
	RoleEditor Role = "editor"
	RoleViewer Role = "viewer"
)

var RoleSet = map[string]bool{
	string(RoleAdmin):  true,
	string(RoleOwner):  true,
	string(RoleEditor): true,
	string(RoleViewer): true,
}

// Resources Definitions
const (
	ResourceAccount Resource = "account"
	ResourcePeriod  Resource = "period"
	ResourceAccess  Resource = "access"
)

// Actions Definitions
const (
	ActionRead   Action = "read"
	ActionWrite  Action = "write"
	ActionDelete Action = "delete"
)

// PermissionSet stores valid "resource:action" pairs
var PermissionSet = map[string]bool{
	string(ResourceAccount) + ":" + string(ActionRead):   true,
	string(ResourceAccount) + ":" + string(ActionWrite):  true,
	string(ResourceAccount) + ":" + string(ActionDelete): true,

	string(ResourcePeriod) + ":" + string(ActionRead):   true,
	string(ResourcePeriod) + ":" + string(ActionWrite):  true,
	string(ResourcePeriod) + ":" + string(ActionDelete): true,

	string(ResourceAccess) + ":" + string(ActionRead):   true,
	string(ResourceAccess) + ":" + string(ActionWrite):  true,
	string(ResourceAccess) + ":" + string(ActionDelete): true,
}
