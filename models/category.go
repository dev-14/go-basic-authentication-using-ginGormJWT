package models

import "time"

type Category struct {
	ID        int    `json:"id,string"`
	Title     string `json:"title" gorm:"unique"`
	CreatedBy int
	User      User `gorm:"foreignKey:CreatedBy" json:"-"`
	// Meta
	CreatedAt time.Time
	UpdatedAt time.Time
}
