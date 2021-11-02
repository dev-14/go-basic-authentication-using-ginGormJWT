package models

type Cart struct {
	CartID int
	UserID int
	User   User `gorm:foreignKey:"UserID"`
	BookID int
	Book   Book `gorm:foreignKey:"BookID"`
	// BookName string
	// Price int
}
