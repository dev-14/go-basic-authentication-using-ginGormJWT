package models

type Cart struct {
	CartID int
	UserID int
	User   User //`gorm:foreignKey:"UserID"`
	BookID int
	Book   Book `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// BookName string
	// Price int
}
