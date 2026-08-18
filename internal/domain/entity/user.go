package entity

type Lang string

const (
	LangEng Lang = "eng"
	LangFa  Lang = "fa"
)

type User struct {
	BaseEntity

	// Relations
	TargetAccounts []TargetAccount `gorm:"foreignKey:OwnerUserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
