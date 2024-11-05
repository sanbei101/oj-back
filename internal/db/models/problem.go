package models

type Problems struct {
	ID          int      `gorm:"column:id"`
	Name        string   `gorm:"column:name"`
	Description string   `gorm:"column:description"`
	Tags        []string `gorm:"column:tags"`
}
