package entity

type Username struct {
	BaseEntity

	UserID          int64
	RoleID          int64
	Username        string
	Email           *string
	IsVerifiedEmail bool
	DayDuration     *int
	Period          *int
	UserLang        Lang
	Description     *string

	// Relations
	User User
	Role Role
}
