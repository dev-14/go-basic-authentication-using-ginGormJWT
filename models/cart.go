package models

type Cart struct {
	CartID int
	UserID int
	User User `gorm:foreignKey:"UserID" json:"-"`
	BookID int
	Book Book `gorm:foreignKey:"BookID" json:"-"`
	// BookName string
	// Price int
}
