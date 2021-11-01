package controllers

import (
	"fmt"
	"gingorm/models"
	"html/template"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// CreateBook godoc
// @Summary CreateBook endpoint is used by the supervisor role user to create a new book.
// @Description CreateBook endpoint is used by the supervisor role user to create a new book
// @Router /api/v1/auth/books/create [post]
// @Tags book
// @Accept json
// @Produce json
// @Param name formData string true "name of the book"
// @Param category_id formData string true "category_id of the book"
func AddToCart(c *gin.Context) {

	//var existingBook models.Book
	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User
	var Book models.Book

	// Check if the current user had admin role.
	if err := models.DB.Where("email = ? AND user_role_id=3", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product can only be added to cart by user"})
		return
	}

	c.Request.ParseForm()

	if c.PostForm("title") == "" {
		ReturnParameterMissingError(c, "title")
		return
	}
	// if c.PostForm("category_id") == "" {
	// 	ReturnParameterMissingError(c, "category_id")
	// 	return
	// }

	title := template.HTMLEscapeString(c.PostForm("title"))
	//fmt.Println(name)
	// category_id := template.HTMLEscapeString(c.PostForm("category_id"))
	// price, err := strconv.Atoi(template.HTMLEscapeString(c.PostForm("price")))
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "can only convert string to int",
	// 	})
	// }

	//Check if the product already exists.
	err := models.DB.Where("title = ?").First(&Book).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "book does not exist."})
		return
	}
	fmt.Println(Book.Price)
	cart := models.Cart{
		// BookName: Book.Title,
		// Price:    Book.Price,
		User: User,
		Book: Book,
	}

	err = models.DB.Create(&cart).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"name":    title,
	})

}
