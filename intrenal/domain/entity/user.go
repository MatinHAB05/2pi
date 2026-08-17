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
	Usernames []Username
	Roles     []Role `gorm:"many2many:user_roles;"`
}
