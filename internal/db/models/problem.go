package models

import (
	"github.com/lib/pq"
)

type Problems struct {
	ID          int            `gorm:"column:id" json:"id,omitempty"`
	Name        string         `gorm:"column:name" json:"name,omitempty"`
	Description string         `gorm:"column:description" json:"description,omitempty"`
	Tags        pq.StringArray `gorm:"column:tags;type:text[]" json:"tags,omitempty"`
}
