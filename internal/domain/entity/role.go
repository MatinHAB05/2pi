package entity

type RoleName uint

const (
	RoleAdmin RoleName = iota + 1
	RoleOwner
	RoleEditor
	RoleViewer
)

var roleNames = map[RoleName]string{
	RoleAdmin:  "admin",
	RoleOwner:  "owner",
	RoleEditor: "editor",
	RoleViewer: "viewer",
}
