package main

import (
	//"net/http"

	"gingorm/controllers"
	"gingorm/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	models.ConnectDataBase()
	r.GET("/books", controllers.FindBooks)
	r.POST("/books", controllers.CreateBook)
	r.GET("/books/:id", controllers.FindBook)
	//r.PATCH("/books/:id", controllers.UpdateBook)
	r.DELETE("/books/:id", controllers.DeleteBook)
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/user", controllers.User)

	r.Run()
}
