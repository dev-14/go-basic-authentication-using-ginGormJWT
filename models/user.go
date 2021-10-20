package models

type UserRole struct {
	Id   int    `json:"id,string"`
	Role string `json:"role"`
}

type User struct {
	ID         uint   `json:"id"`
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Email      string `json:"email" gorm:"unique"`
	UserRoleID int    `json:"role_id,string"`
	Password   string `json:"-"`
}
