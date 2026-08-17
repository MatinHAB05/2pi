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
