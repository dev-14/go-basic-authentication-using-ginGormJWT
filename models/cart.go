package models

type Cart struct {
	CartID int `gorm:"autoIncrement"`
	UserID int
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	BookID int
	Book   Book `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// BookName string
	// Price int
}
