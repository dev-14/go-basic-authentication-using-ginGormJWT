package controllers

import (
	"fmt"
	"gingorm/models"
	"html/template"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// CreateCategory godoc
// @Summary CreateCategory endpoint is used by admin to create category.
// @Description API Endpoint to register the user with the role of Supervisor or Admin.
// @Router /api/v1/auth/category/create [post]
// @Tags category
// @Accept json
// @Produce json
// @Param name formData string true "name of the category"
func CreateCategory(c *gin.Context) {
	var existingCategory models.Category
	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User

	// Check if the current user had admin role.
	if err := models.DB.Where("email = ? AND user_role_id=1", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category can only be added by admin user"})
	}

	c.Request.ParseForm()

	if c.PostForm("name") == "" {
		ReturnParameterMissingError(c, "name")
	}

	category_title := template.HTMLEscapeString(c.PostForm("name"))
	fmt.Println(category_title)
	fmt.Println("category printed")
	// Check if the category already exists.

	err := models.DB.Where("title = ?", category_title).First(&existingCategory).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category already exists."})
		return
	}

	cat := models.Category{
		CategoryName: category_title,
		CreatedBy:    User.ID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = models.DB.Create(&cat).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":   cat.ID,
		"name": cat.CategoryName,
	})

}

type ReturnedCategory struct {
	ID           int    `json:"id,string"`
	CategoryName string `json:"name"`
}

// ListAllCategories godoc
// @Summary ListAllCategories endpoint is used to list all the categories.
// @Description API Endpoint to register the user with the role of Supervisor or Admin.
// @Router /api/v1/auth/category/ [get]
// @Tags category
// @Accept json
// @Produce json
func ListAllCategories(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User
	var Categories []models.Category
	var ExistingCategories []ReturnedCategory

	if err := models.DB.Where("email = ?", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	models.DB.Model(Categories).Find(&ExistingCategories)
	c.JSON(http.StatusOK, ExistingCategories)
	return
}

// GetCategory godoc
// @Summary GetCategory endpoint is used to get info of a category..
// @Description GetCategory endpoint is used to get info of a category.
// @Router /api/v1/auth/category/:id/ [get]
// @Tags product
// @Accept json
// @Produce json
func GetCategory(c *gin.Context) {
	var existingCategory models.Category

	// Check if the category already exists.
	err := models.DB.Where("id = ?", c.Param("id")).First(&existingCategory).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category doesnot exists."})
		return
	}

	// GET FROM CACHE FIRST
	c.JSON(http.StatusOK, gin.H{"category": existingCategory})
	return
}

// GetCategory godoc
// @Summary UpdateCatagory endpoint is used to get info of a category..
// @Description UpdateCatagory endpoint is used to get info of a category.
// @Router /api/v1/auth/category/:id/ [PUT]
// @Tags product
// @Accept json
// @Produce json
func UpdateCategory(c *gin.Context) {
	var existingCategory models.Category
	var UpdateCategory models.Category

	// Check if the product already exists.
	err := models.DB.Where("id = ?", c.Param("id")).First(&existingCategory).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category doesnot exists."})
		return
	}

	if err := c.ShouldBindJSON(&UpdateCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Model(&existingCategory).Updates(UpdateCategory)
}
