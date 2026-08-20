package entity

type Lang string

const (
	LangEng Lang = "eng"
	LangFa  Lang = "fa"
)

type User struct {
	BaseEntity

	FirstName string
	LastName  string
	Username  string
	UserLang  Lang `gorm:"type:lang;not null;default:'fa'"`

	// Relations
	TargetAccounts []TargetAccount `gorm:"foreignKey:OwnerUserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
