package models

import (
	"github.com/lib/pq"
)

type Problems struct {
	ID          int            `gorm:"column:id"`
	Name        string         `gorm:"column:name"`
	Description string         `gorm:"column:description"`
	Tags        pq.StringArray `gorm:"column:tags;type:text[]"`
}
