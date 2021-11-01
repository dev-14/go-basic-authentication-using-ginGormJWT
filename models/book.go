package models

import "time"

type Book struct {
	ID         int    `json:"id,string"`
	Title      string `json:"title" gorm:"unique"`
	CreatedBy  int
	User       User `gorm:"foreignKey:CreatedBy" json:"-"`
	CategoryId int  `json:"category_id,string"`
	Price      int  `json:"price"`
	// Meta
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type BookImage struct {
	ID     int    `json:"id,string"`
	URL    string `json:"uri"`
	BookId int
	Book   Book `gorm:"foreignKey:BookId"`

	// META
	CreatedAt time.Time
}
