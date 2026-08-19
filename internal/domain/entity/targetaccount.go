package entity

type TargetAccount struct {
	BaseEntity

	OwnerUserID int64 `gorm:"not null;index"`
	DayDuration int
	Period      int
	Description string
	Enable      bool
	// Relations
	Owner User `gorm:"foreignKey:OwnerUserID;references:ID"`
}

// Email           *string `gorm:"size:255"`
// IsVerifiedEmail bool    `gorm:"default:false"`
