package controllers

import (
	"fmt"
	"gingorm/models"
	"html/template"
	"net/http"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const SecretKey = "secret"

var Flag string

type tempUser struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	// Email           string `json:"email" binding:""`
	// Mobile          string `json:"mobile" binding:""`
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmpassword" binding:"required"`
}

type RedisCache struct {
	Id     int
	Email  string
	RoleId int
}

func ReturnParameterMissingError(c *gin.Context, parameter string) {
	var err = fmt.Sprintf("Required parameter %s missing.", parameter)
	c.JSON(http.StatusBadRequest, gin.H{"error": err})
}

// @Summary register endpoint is used for customer registration. ( Supervisors/admin can be added only by admin. )
// @Description API Endpoint to register the user with the role of customer.
// @Router /api/v1/register [post]
// @Tags auth
// @Accept json
// @Produce json
// @Success 200
// @Param email formData string true "Email of the user"
// @Param first_name formData string true "First name of the user"
// @Param last_name formData string true "Last name of the user"
// @Param password formData string true "Password of the user"
// @Param confirm_password formData string true "Confirm password."
func Register(c *gin.Context) {
	var tempUser tempUser
	var Role models.UserRole

	c.Request.ParseForm()
	paramList := []string{"username", "first_name", "last_name", "password", "confirmpassword"}

	for _, param := range paramList {
		if c.PostForm(param) == "" {
			ReturnParameterMissingError(c, param)
		}
	}

	// if err := c.ShouldBindJSON(&tempUser); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	tempUser.Username = template.HTMLEscapeString(c.PostForm("username"))
	tempUser.FirstName = template.HTMLEscapeString(c.PostForm("first_name"))
	tempUser.LastName = template.HTMLEscapeString(c.PostForm("last_name"))
	tempUser.Password = template.HTMLEscapeString(c.PostForm("password"))
	tempUser.ConfirmPassword = template.HTMLEscapeString(c.PostForm("confirmpassword"))

	if tempUser.Password != tempUser.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both passwords do not match."})
		return
	}

	ispasswordstrong, _ := IsPasswordStrong(tempUser.Password)
	if !ispasswordstrong {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong."})
		return
	}

	// Check if the user already exists.
	if DoesUserExist(tempUser.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists."})
		return
	}

	encryptedPassword, error := HashPassword(tempUser.Password)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
		return
	}

	err := models.DB.Where("role= ?", "customer").First(&Role).Error
	if err != nil {
		fmt.Println("err ", err.Error())
		return
	}

	flag := strings.Index(tempUser.Username, "@")
	if flag == -1 {
		fmt.Println(len([]rune(tempUser.Username)))
		if len([]rune(tempUser.Username)) != 10 {
			c.JSON(404, gin.H{
				"error": "Invalid Mobile Number",
			})
			return
		}
		// Flag = "mobile"
		SanitizedUser := models.User{
			FirstName:  tempUser.FirstName,
			LastName:   tempUser.LastName,
			Mobile:     tempUser.Username,
			Password:   encryptedPassword,
			UserRoleID: Role.Id, //This endpoint will be used only for customer registration.
			CreatedAt:  time.Now(),
			IsActive:   true,
		}
		errs := models.DB.Create(&SanitizedUser).Error
		if errs != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured while input to db"})
			return
		}
	} else {
		// Flag = "email"
		SanitizedUser := models.User{
			FirstName:  tempUser.FirstName,
			LastName:   tempUser.LastName,
			Email:      tempUser.Username,
			Password:   encryptedPassword,
			UserRoleID: Role.Id, //This endpoint will be used only for customer registration.
			CreatedAt:  time.Now(),
			IsActive:   true,
		}
		errs := models.DB.Create(&SanitizedUser).Error
		if errs != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured while input to db"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"msg": "User created successfully"})

}

type login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// redisClient := Redis.createclient()

// Login godoc
// @Summary Login endpoint is used by the user to login.
// @Description API Endpoint to register the user with the role of customer.
// @Router /api/v1/login [post]
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 description {object}
// @Param email formData string true "email id"
// @Param password formData string true "password"
func Login(c *gin.Context) (interface{}, error) {
	// var loginVals login
	// var User User
	var user models.User
	var count int64

	username := template.HTMLEscapeString(c.PostForm("username"))
	password := template.HTMLEscapeString(c.PostForm("password"))

	flag := strings.Index(username, "@")

	if flag == -1 {
		Flag = "mobile"
		fmt.Println("mobile hai")

	} else {
		Flag = "email"
		fmt.Println("email hai")
	}
	// var user models.User
	// if err := c.ShouldBind(&loginVals); err != nil {
	// 	return "", jwt.ErrMissingLoginValues
	// }
	fmt.Println(Flag)
	// email := loginVals.Email
	// First check if the user exist or not...
	if Flag == "email" {
		models.DB.Where("email = ?", username).First(&user).Count(&count)
		// if count == 0 {
		// 	return nil, jwt.ErrFailedAuthentication
		// }
	} else if Flag == "mobile" {
		models.DB.Where("mobile = ?", username).First(&user).Count(&count)
		if count == 0 {
			return nil, jwt.ErrFailedAuthentication
		}
	}
	if CheckCredentials(username, password, models.DB) {
		NewRedisCache(c, user)
		if Flag == "email" {
			return &models.User{
				Email: username,
			}, nil
		} else {
			return &models.User{
				Mobile: username,
			}, nil
		}

	}
	// fmt.Println("set value ", loginVals.Email)
	// err := rdb.Set("email", loginVals.Email, 0).Err()
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{
	// 		"error": "error in redis",
	// 	})
	// }

	return nil, jwt.ErrFailedAuthentication
}

func CheckUserLevel(c *gin.Context) {
	if IsAdmin(c) {
		c.JSON(200, gin.H{
			"status": "admin",
		})
	} else if IsSupervisor(c) {
		c.JSON(200, gin.H{
			"status": "supervisor",
		})
	} else {
		c.JSON(200, gin.H{
			"status": "user"})
	}

}

// CreateSupervisor godoc
// @Summary CreateSupervisor endpoint is used by the admin role user to create a new admin or supervisor account.
// @Description API Endpoint to register the user with the role of Supervisor or Admin.
// @Router /api/v1/auth/supervisor/create [post]
// @Tags supervisor
// @Accept json
// @Produce json
// @Param login formData tempUser true "Info of the user"
func CreateSupervisor(c *gin.Context) {
	//fmt.Println("supervisor api hit")

	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}

	// Create a user with the role of supervisor.
	var tempUser tempUser
	var Role models.UserRole

	c.Request.ParseForm()
	paramList := []string{"first_name", "last_name", "email", "password", "confirmpassword"}

	for _, param := range paramList {
		if c.PostForm(param) == "" {
			ReturnParameterMissingError(c, param)
		}
	}

	tempUser.Username = template.HTMLEscapeString(c.PostForm("username"))
	tempUser.FirstName = template.HTMLEscapeString(c.PostForm("first_name"))
	tempUser.LastName = template.HTMLEscapeString(c.PostForm("last_name"))
	tempUser.Password = template.HTMLEscapeString(c.PostForm("password"))
	tempUser.ConfirmPassword = template.HTMLEscapeString(c.PostForm("confirmpassword"))

	//check if the password is strong and matches the password policy
	//length > 8, atleast 1 upper case, atleast 1 lower case, atleast 1 symbol
	// ispasswordstrong, _ := IsPasswordStrong(tempUser.Password)
	// if ispasswordstrong == false {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong."})
	// 	return
	// }

	if tempUser.Password != tempUser.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both passwords do not match."})
	}

	ispasswordstrong, _ := IsPasswordStrong(tempUser.Password)
	if !ispasswordstrong {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong."})
		return
	}

	// Check if the user already exists.
	if DoesUserExist(tempUser.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists."})
		return
	}

	encryptedPassword, error := HashPassword(tempUser.Password)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
		return
	}

	err := models.DB.Where("role= ?", "supervisor").First(&Role).Error
	if err != nil {
		fmt.Println("err ", err.Error())
		return
	}

	flag := strings.Index(tempUser.Username, "@")
	if flag == -1 {
		if len(tempUser.Username) != 10 {
			c.JSON(404, gin.H{
				"error": "Invalid Mobile Number",
			})
			return
		}

		SanitizedUser := models.User{
			FirstName:  tempUser.FirstName,
			LastName:   tempUser.LastName,
			Mobile:     tempUser.Username,
			Password:   encryptedPassword,
			UserRoleID: Role.Id, //This endpoint will be used only for customer registration.
			CreatedAt:  time.Now(),
			IsActive:   true,
		}
		errs := models.DB.Create(&SanitizedUser).Error
		if errs != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured while input to db"})
			return
		}
	} else {

		SanitizedUser := models.User{
			FirstName:  tempUser.FirstName,
			LastName:   tempUser.LastName,
			Email:      tempUser.Username,
			Password:   encryptedPassword,
			UserRoleID: Role.Id, //This endpoint will be used only for customer registration.
			CreatedAt:  time.Now(),
			IsActive:   true,
		}
		errs := models.DB.Create(&SanitizedUser).Error
		if errs != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured while input to db"})
			return
		}
	}

	// SanitizedUser := models.User{
	// 	FirstName:  tempUser.FirstName,
	// 	LastName:   tempUser.LastName,
	// 	Email:      tempUser.Email,
	// 	Password:   encryptedPassword,
	// 	UserRoleID: Role.Id, //This endpoint will be used only for customer registeration.
	// 	CreatedAt:  time.Now(),
	// 	IsActive:   true,
	// }

	// errs := models.DB.Create(&SanitizedUser).Error
	// if errs != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
	// 	return
	// }
	c.JSON(http.StatusCreated, gin.H{"msg": "User created successfully"})

}

// CreateAdmin godoc
// @Summary CreateAdmin endpoint is used by the admin role user to create a new admin or supervisor account.
// @Description API Endpoint to register the user with the role of Supervisor or Admin.
// @Router /api/v1/auth/admin/create [post]
// @Tags admin
// @Accept json
// @Produce json
// @Param login formData tempUser true "Info of the user"
func CreateAdmin(c *gin.Context) {
	//var User models.User
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	// user_email := rdb.
	//fmt.Println("test line")
	//fmt.Println(user_email)
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{
	// 		"error": "redis get not working",
	// 	})
	// }
	// //fmt.Println()

	// if err := models.DB.Where("email = ? AND user_role_id=1", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }

	var tempUser tempUser
	var Role models.UserRole

	c.Request.ParseForm()
	paramList := []string{"first_name", "last_name", "email", "password", "confirmpassword"}

	for _, param := range paramList {
		if c.PostForm(param) == "" {
			ReturnParameterMissingError(c, param)
			return
		}
	}

	tempUser.Username = template.HTMLEscapeString(c.PostForm("username"))
	tempUser.FirstName = template.HTMLEscapeString(c.PostForm("first_name"))
	tempUser.LastName = template.HTMLEscapeString(c.PostForm("last_name"))
	tempUser.Password = template.HTMLEscapeString(c.PostForm("password"))
	tempUser.ConfirmPassword = template.HTMLEscapeString(c.PostForm("confirmpassword"))

	// fmt.Println("debug start")
	fmt.Println(tempUser.FirstName, tempUser.LastName, tempUser.Password)
	// fmt.Println("debug end")

	if tempUser.Password != tempUser.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both passwords do not match."})
		return
	}

	// check if the password is strong and matches the password policy
	// length > 8, atleast 1 upper case, atleast 1 lower case, atleast 1 symbol
	// ispasswordstrong, _ := IsPasswordStrong(tempUser.Password)
	// if ispasswordstrong == false {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong."})
	// 	return
	// }

	ispasswordstrong, _ := IsPasswordStrong(tempUser.Password)
	if !ispasswordstrong {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong."})
		return
	}

	// Check if the user already exists.
	if DoesUserExist(tempUser.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists."})
		return
	}

	encryptedPassword, error := HashPassword(tempUser.Password)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
		return
	}

	err := models.DB.Where("role= ?", "admin").First(&Role).Error
	if err != nil {
		fmt.Println("error ", err.Error())
		return
	}

	flag := strings.Index(tempUser.Username, "@")
	if flag == -1 {
		if len(tempUser.Username) != 10 {
			c.JSON(404, gin.H{
				"error": "Invalid Mobile Number",
			})
			return
		}
		SanitizedUser := models.User{
			FirstName:  tempUser.FirstName,
			LastName:   tempUser.LastName,
			Mobile:     tempUser.Username,
			Password:   encryptedPassword,
			UserRoleID: Role.Id, //This endpoint will be used only for customer registration.
			CreatedAt:  time.Now(),
			IsActive:   true,
		}
		errs := models.DB.Create(&SanitizedUser).Error
		if errs != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured while input to db"})
			return
		}
	} else {
		SanitizedUser := models.User{
			FirstName:  tempUser.FirstName,
			LastName:   tempUser.LastName,
			Email:      tempUser.Username,
			Password:   encryptedPassword,
			UserRoleID: Role.Id, //This endpoint will be used only for customer registration.
			CreatedAt:  time.Now(),
			IsActive:   true,
		}
		errs := models.DB.Create(&SanitizedUser).Error
		if errs != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured while input to db"})
			return
		}
	}

	// SanitizedUser := models.User{
	// 	FirstName:  tempUser.FirstName,
	// 	LastName:   tempUser.LastName,
	// 	Email:      tempUser.Email,
	// 	Password:   encryptedPassword,
	// 	UserRoleID: Role.Id, //This endpoint will be used only for customer registeration.
	// 	CreatedAt:  time.Now(),
	// 	IsActive:   true,
	// }

	// errs := models.DB.Create(&SanitizedUser).Error
	// if errs != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
	// 	return
	// }
	c.JSON(http.StatusCreated, gin.H{"msg": "User created successfully"})

}

func MyProfile(c *gin.Context) {
	var User models.User

	username, _ := models.Rdb.HGet("user", "username").Result()

	if username == "" {
		fmt.Println("Redis empty....checking Database for user...")
		err := FillRedis(c)
		if err != nil {
			c.JSON(404, gin.H{
				"error": "something went wrong with redis",
			})
			return
		}
	}
	username, _ = models.Rdb.HGet("user", "username").Result()

	if Flag == "email" {
		if err := models.DB.Where("email = ?", username).First(&User).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
	} else {
		if err := models.DB.Where("mobile = ?", username).First(&User).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
	}
	c.JSON(200, &User)

}

// type updateUser struct {
// 	Field string `json:"field"`
// }

func UpdateUser(c *gin.Context) {
	var user models.User
	var existingUser models.User
	var updateUser models.User
	var count int64

	id, _ := models.Rdb.HGet("user", "ID").Result()
	_ = id

	err := models.DB.Where("id = ?", c.Param("id")).First(&existingUser).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user doesnot exists."})
		return
	}

	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if existingUser.Mobile == "" {
		fmt.Println("mobile nahi hai")
		models.DB.Where("mobile = ?", updateUser.Mobile).First(&user).Count(&count)
		if count != 0 {
			c.JSON(404, gin.H{"error": "mobile number linked with another user, pls try different mobile number"})
			return
		}

	} else if existingUser.Email == "" {
		fmt.Println("email nahi hai")
		models.DB.Where("email = ?", updateUser.Email).First(&user).Count(&count)
		if count != 0 {
			c.JSON(404, gin.H{"error": "email address linked with another user, pls try different email"})
			return
		}
	}

	models.DB.Model(&existingUser).Updates(updateUser)
}

func GetAllUsers(c *gin.Context) {
	// var User []models.User
	var existingUsers []models.User

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
		c.JSON(404, gin.H{
			"error": "unauthorized",
		})
		return
	}
	// models.DB.Model(User).Find(existingUsers)
	result := models.DB.Find(&existingUsers)
	fmt.Println(result)
	c.JSON(200, existingUsers)
}
