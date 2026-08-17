package entity

type Lang string

const (
	LangEng Lang = "eng"
	LangFa  Lang = "fa"
)

type User struct {
	//TODO : CHECK `ID        int64 `gorm:"primaryKey;autoIncrement:false"` // Base On Telegram User ID
	BaseEntity

	// Relations
	TargetAccounts []TargetAccount `gorm:"foreignKey:OwnerUserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
