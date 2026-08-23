package entity

type Article struct {
	BaseEntity
	Title       string
	Description string
	ImageURL    string
	URL         string
}
