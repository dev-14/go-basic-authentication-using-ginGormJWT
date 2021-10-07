package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var DB *gorm.DB

func ConnectDataBase() {
	database := fmt.Sprintf("host=localhost port=5432 user=postgres dbname=postgres password=admin sslmode=disable")
	fmt.Println("conname is\t", database)
	connection, err := gorm.Open("postgres", database)

	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = connection
	connection.AutoMigrate(&Book{})
	connection.AutoMigrate(&User{})

}
