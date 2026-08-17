package entity

type RoleName string

const (
	RoleAdmin  RoleName = "admin"
	RoleOwner  RoleName = "owner"
	RoleEditor RoleName = "editor"
	RoleViewer RoleName = "viewer"
)

type Role struct {
	BaseEntity

	Name RoleName

	// Relations
	Users     []User `gorm:"many2many:user_roles;"`
	Usernames []Username
}
