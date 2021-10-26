package controllers

import (
	"gingorm/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateBook godoc
// @Summary CreateBook endpoint is used by the supervisor role user to create a new book.
// @Description CreateBook endpoint is used by the supervisor role user to create a new book
// @Router /api/v1/auth/product/create [post]
// @Tags product
// @Accept json
// @Produce json
func CreateBook(c *gin.Context) {

	var existingBook models.Book
	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User
	var category models.Category

	// Check if the current user had admin role.
	if err := models.DB.Where("email = ? AND user_role_id=2", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product can only be added by supervisor user"})
		return
	}

	c.Request.ParseForm()

	if c.PostForm("name") == "" {
		ReturnParameterMissingError(c, "name")
		return
	}
	if c.PostForm("category_id") == "" {
		ReturnParameterMissingError(c, "category_id")
		return
	}

	title := template.HTMLEscapeString(c.PostForm("name"))
	category_id := template.HTMLEscapeString(c.PostForm("category_id"))
	price, err := strconv.Atoi(template.HTMLEscapeString(c.PostForm("price")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "can only convert string to int",
		})
	}

	//Check if the product already exists.
	err = models.DB.Where("title = ?").First(&existingBook).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product already exists."})
		return
	}

	// Check if the category exists
	err = models.DB.First(&category, category_id).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category does not exists."})
		return
	}

	book := models.Book{
		Title:      title,
		CategoryId: category.ID,
		Price:      price,
		CreatedBy:  User.ID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = models.DB.Create(&book).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":   book.ID,
		"name": book.Title,
	})

}

// UpdateBook godoc
// @Summary UpdateBook endpoint is used by the supervisor role user to update a new book.
// @Description Updatebook endpoint is used by the supervisor role user to update a new book
// @Router /api/v1/auth/product/:id/ [PATCH]
// @Tags book
// @Accept json
// @Produce json
func UpdateBook(c *gin.Context) {
	var existingBook models.Book
	var updateBook models.Book
	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User

	// Check if the current user had admin role.
	if err := models.DB.Where("email = ? AND user_role_id=2", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product can only be updated by supervisor user"})
		return
	}

	// Check if the product already exists.
	err := models.DB.Where("id = ?", c.Param("id")).First(&existingBook).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product doesnot exists."})
		return
	}

	if err := c.ShouldBindJSON(&updateBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Model(&existingBook).Updates(updateBook)

}

type ReturnedBook struct {
	ID         int    `json:"id,string"`
	Title      string `json:"name"`
	CategoryId int    `json:"category_id"`
}

// GetBook godoc
// @Summary GetBook endpoint is used to get info of a book..
// @Description GetBook endpoint is used to get info of a book.
// @Router /api/v1/auth/product/:id/ [get]
// @Tags product
// @Accept json
// @Produce json
func GetBook(c *gin.Context) {
	var existingBook models.Book

	// Check if the product already exists.
	err := models.DB.Where("id = ?", c.Param("id")).First(&existingBook).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product doesnot exists."})
		return
	}

	// GET FROM CACHE FIRST
	c.JSON(http.StatusOK, gin.H{"product": existingBook})
	return
}

// ListAllBook godoc
// @Summary ListAllBook endpoint is used to list all book.
// @Description API Endpoint to register the user with the role of Supervisor or Admin.
// @Router /api/v1/auth/book/ [get]
// @Tags book
// @Accept json
// @Produce json
func ListAllProduct(c *gin.Context) {

	// allProduct := []models.Product{}
	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User
	var Product []models.Book
	var existingBook []ReturnedBook

	if err := models.DB.Where("email = ?", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	models.DB.Model(Product).Find(&existingBook)
	c.JSON(http.StatusOK, existingBook)
	return
}

func generateFilePath(id string, extension string) string {
	// Generate random file name for the new uploaded file so it doesn't override the old file with same name
	newFileName := uuid.New().String() + extension

	projectFolder, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	localS3Folder := projectFolder + "/locals3/"
	productImageFolder := localS3Folder + id + "/"

	if _, err := os.Stat(productImageFolder); os.IsNotExist(err) {
		os.Mkdir(productImageFolder, os.ModeDir)
	}

	imagePath := productImageFolder + newFileName
	return imagePath
}
