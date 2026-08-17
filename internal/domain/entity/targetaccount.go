package entity

type TargetAccount struct {
	BaseEntity

	OwnerUserID     int64   `gorm:"not null;index"`
	Username        string  `gorm:"size:100;not null;unique"`
	Email           *string `gorm:"size:255"`
	IsVerifiedEmail bool    `gorm:"default:false"`
	DayDuration     *int
	Period          *int
	UserLang        Lang `gorm:"type:lang;not null;default:'fa'"`
	Description     *string

	// Relations
	Owner User `gorm:"foreignKey:OwnerUserID;references:ID"`
}
