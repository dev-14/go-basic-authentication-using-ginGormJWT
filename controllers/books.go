package controllers

import (
	"gingorm/models"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func CreateBook(c *gin.Context) {

	var existingProduct models.Book
	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User
	var category models.Category

	// Check if the current user had admin role.
	if err := models.DB.Where("email = ? AND user_role_id=2", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product can only be added by supervisor user"})
		return
	}

	// c.Request.ParseForm()

	// if c.PostForm("name") == "" {
	// 	ReturnParameterMissingError(c, "name")
	// 	return
	// }
	// if c.PostForm("category_id") == "" {
	// 	ReturnParameterMissingError(c, "category_id")
	// 	return
	// }

	// product_title := template.HTMLEscapeString(c.PostForm("name"))
	// category_id := template.HTMLEscapeString(c.PostForm("category_id"))

	type book struct {
		Title      string
		CategoryId int
		Price      int
	}
	var newbook book

	if err := c.ShouldBindJSON(&newbook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
	}

	// Check if the product already exists.
	err := models.DB.Where("title = ?", newbook.Title).First(&existingProduct).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product already exists."})
		return
	}

	// Check if the category exists
	err = models.DB.First(&category, newbook.CategoryId).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category does not exists."})
		return
	}
	createdBook := models.Book{
		Title:      newbook.Title,
		CategoryId: newbook.CategoryId,
		Price:      newbook.Price,
		CreatedBy:  User.ID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = models.DB.Create(&createdBook).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":   createdBook.ID,
		"name": createdBook.Title,
	})

}
