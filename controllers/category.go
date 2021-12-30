package controllers

import (
	"fmt"
	"gingorm/models"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/go-redis/redis/v8"
)

// var Rdb = redis.NewClient(&redis.Options{
// 	Addr:     "localhost:6379",
// 	Password: "", // no password set
// 	DB:       0,  // use default DB
// })

// CreateCategory godoc
// @Summary CreateCategory endpoint is used by admin to create category.
// @Description CreateCategory endpoint is used by admin to create category.
// @Router /api/v1/auth/category/create [post]
// @Tags category
// @Accept json
// @Produce json
// @Param name formData string true "name of the category"
func CreateCategory(c *gin.Context) {
	var existingCategory models.Category
	// claims := jwt.ExtractClaims(c)
	// user_email, _ := claims["email"]
	//var User models.User
	email := c.GetString("user_email")
	fmt.Println(models.Rdb.HGetAll(email))
	// user_email, err := Rdb.HGet("user", "email").Result()
	id, _ := models.Rdb.HGet("user", "ID").Result()
	ID, _ := strconv.Atoi(id)
	roleId, _ := models.Rdb.HGet("user", "RoleID").Result()

	if roleId == "" {
		fmt.Println("Redis empty....checking Database for user...")
		err := FillRedis(c)
		if err != nil {
			c.JSON(404, gin.H{
				"error": "something went wrong with redis",
			})
			return
		}
	}

	roleId, _ = models.Rdb.HGet("user", "RoleID").Result()

	if roleId != "1" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category can only be updated by admin user"})
		return
	}

	// // Check if the current user had admin role.
	// if err := models.DB.Where("email = ? AND user_role_id=1", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Category can only be added by admin user"})
	// }

	c.Request.ParseForm()
	var flag bool
	if c.PostForm("name") == "" {
		ReturnParameterMissingError(c, "name")
		flag = true
	}
	category_title := template.HTMLEscapeString(c.PostForm("name"))
	// fmt.Println(category_title)
	// fmt.Println("category printed")
	// Check if the category already exists.
	if flag == true {
		return
	}
	err := models.DB.Where("title = ?", category_title).First(&existingCategory).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category already exists."})
		return
	}

	cat := models.Category{
		CategoryName: category_title,
		CreatedBy:    ID,
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
// @Description ListAllCategories endpoint is used to list all the categories.
// @Router /api/v1/auth/category/ [get]
// @Tags category
// @Accept json
// @Produce json
func ListAllCategories(c *gin.Context) {

	// claims := jwt.ExtractClaims(c)
	// user_email, _ := claims["email"]
	// var User models.User
	var Categories []models.Category
	var ExistingCategories []ReturnedCategory
	//email := c.GetString("user_email")
	username, _ := models.Rdb.HGet("user", "username").Result()

	// if err := models.DB.Where("email = ?", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }

	if !IsAuthorized(username) {
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
// @Tags category
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

// UpdateCategory godoc
// @Summary UpdateCategory endpoint is used to get info of a category..
// @Description UpdateCategory endpoint is used to get info of a category.
// @Router /api/v1/auth/category/:id/ [PUT]
// @Tags category
// @Accept json
// @Produce json
func UpdateCategory(c *gin.Context) {
	// claims := jwt.ExtractClaims(c)
	// user_email, _ := claims["email"]
	//var User models.User
	var existingCategory models.Category
	var UpdateCategory models.Category
	//email := c.GetString("user_email")
	//user_email, _ := Rdb.HGet("user", "email").Result()
	id, _ := models.Rdb.HGet("user", "RoleID").Result()

	if id == "" {
		fmt.Println("Redis empty....checking Database for user...")
		err := FillRedis(c)
		if err != nil {
			c.JSON(404, gin.H{
				"error": "something went wrong with redis",
			})
			return
		}
	}

	id, _ = models.Rdb.HGet("user", "RoleID").Result()

	if id != "1" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category can only be updated by admin user"})
		return
	}

	// if err := models.DB.Where("email = ? AND user_role_id=1", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Category can only be updated by admin user"})
	// 	return
	// }
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
