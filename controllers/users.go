package controllers

import (
	"fmt"
	"gingorm/models"
	"html/template"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	//"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const SecretKey = "secret"

// @Summary register endpoint is used for customer registration. ( Supervisors/admin can be added only by admin. )
// @Description API Endpoint to register the user with the role of customer.
// @Router /api/v1/register [post]
// @Tags auth
// @Accept json
// @Produce json
type tempUser struct {
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmpassword" binding:"required"`
}

func ReturnParameterMissingError(c *gin.Context, parameter string) {
	var err = fmt.Sprintf("Required parameter %s missing.", parameter)
	c.JSON(http.StatusBadRequest, gin.H{"error": err})
}

func Register(c *gin.Context) {
	var tempUser tempUser
	var Role models.UserRole

	c.Request.ParseForm()
	paramList := []string{"email", "first_name", "last_name", "password", "confirmpassword"}

	for _, param := range paramList {
		if c.PostForm(param) == "" {
			ReturnParameterMissingError(c, param)
		}
	}

	// if err := c.ShouldBindJSON(&tempUser); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	tempUser.Email = template.HTMLEscapeString(c.PostForm("email"))
	tempUser.FirstName = template.HTMLEscapeString(c.PostForm("first_name"))
	tempUser.LastName = template.HTMLEscapeString(c.PostForm("last_name"))
	tempUser.Password = template.HTMLEscapeString(c.PostForm("password"))
	tempUser.ConfirmPassword = template.HTMLEscapeString(c.PostForm("confirmpassword"))

	if tempUser.Password != tempUser.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both passwords do not match."})
	}

	ispasswordstrong, _ := IsPasswordStrong(tempUser.Password)
	if ispasswordstrong == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong."})
		return
	}

	// Check if the user already exists.
	if DoesUserExist(tempUser.Email) {
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

	SanitizedUser := models.User{
		FirstName:  tempUser.FirstName,
		LastName:   tempUser.LastName,
		Email:      tempUser.Email,
		Password:   encryptedPassword,
		UserRoleID: Role.Id, //This endpoint will be used only for customer registration.
		CreatedAt:  time.Now(),
		IsActive:   true,
	}

	errs := models.DB.Create(&SanitizedUser).Error
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"msg": "User created successfully"})

}

type login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Signin godoc
// @Summary Login endpoint is used by the user to login.
// @Description API Endpoint to register the user with the role of customer.
// @Router /api/v1/login [post]
// @Tags auth
// @Accept json
// @Produce json
func Login(c *gin.Context) (interface{}, error) {
	var loginVals login
	// var User User
	var users []models.User
	var count int64
	// var user models.User
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	email := loginVals.Email
	// First check if the user exist or not...
	models.DB.Where("email = ?", email).Find(&users).Count(&count)
	if count == 0 {
		return nil, jwt.ErrFailedAuthentication
	}
	if CheckCredentials(loginVals.Email, loginVals.Password, models.DB) == true {
		return &models.User{
			Email: email,
		}, nil
	}
	return nil, jwt.ErrFailedAuthentication
}

// func Login(c *gin.Context) {
// 	var data map[string]string

// 	err := c.ShouldBindJSON(&data)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error": err,
// 		})
// 	}

// 	var user models.User

// 	models.DB.Where("email=?", data["email"]).First(&user)
// 	if user.ID == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"message": "user not found",
// 		})
// 	} else if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "wrong password",
// 		})
// 	} else {
// 		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
// 			Issuer:    strconv.Itoa(int(user.ID)),
// 			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
// 		})

// 		token, err := claims.SignedString([]byte(SecretKey))

// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"message": "could not login",
// 			})
// 		}
// 		// _ = token

// 		// c.SetCookie("jwt", token, 60*60*24, "", "", true, true)
// 		// http.SetCookie(c.Writer, &http.Cookie{
// 		// 	Name:    user.FirstName,
// 		// 	Value:   token,
// 		// 	Expires: time.Now().Add(time.Hour * 24),
// 		// })
// 		c.Header("token", token)

// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "successfully logged in",
// 			"token":   token,
// 		})
// 	}

// }

// CreateSupervisorOrAdmin godoc
// @Summary CreateSupervisor endpoint is used by the admin role user to create a new admin or supervisor account.
// @Description API Endpoint to register the user with the role of Supervisor or Admin.
// @Router /api/v1/auth/supervisor/create [post]
// @Tags supervisor
// @Accept json
// @Produce json
// @Param login formData TempUser true "Info of the user"
func CreateSupervisor(c *gin.Context) {
	fmt.Println("supervisor api hit")
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

	tempUser.Email = template.HTMLEscapeString(c.PostForm("email"))
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

	// Check if the user already exists.
	if DoesUserExist(tempUser.Email) {
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

	SanitizedUser := models.User{
		FirstName:  tempUser.FirstName,
		LastName:   tempUser.LastName,
		Email:      tempUser.Email,
		Password:   encryptedPassword,
		UserRoleID: Role.Id, //This endpoint will be used only for customer registeration.
		CreatedAt:  time.Now(),
		IsActive:   true,
	}

	errs := models.DB.Create(&SanitizedUser).Error
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"msg": "User created successfully"})

	return
}

// CreateAdmin godoc
// @Summary CreateAdmin endpoint is used by the admin role user to create a new admin or supervisor account.
// @Description API Endpoint to register the user with the role of Supervisor or Admin.
// @Router /api/v1/auth/admin/create [post]
// @Tags admin
// @Accept json
// @Produce json
// @Param login formData TempUser true "Info of the user"
func CreateAdmin(c *gin.Context) {

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

	tempUser.Email = template.HTMLEscapeString(c.PostForm("email"))
	tempUser.FirstName = template.HTMLEscapeString(c.PostForm("first_name"))
	tempUser.LastName = template.HTMLEscapeString(c.PostForm("last_name"))
	tempUser.Password = template.HTMLEscapeString(c.PostForm("password"))
	tempUser.ConfirmPassword = template.HTMLEscapeString(c.PostForm("confirmpassword"))

	fmt.Println("debug start")
	fmt.Println(tempUser.FirstName, tempUser.LastName, tempUser.Password)
	fmt.Println("debug end")

	if tempUser.Password != tempUser.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both passwords do not match."})
	}

	// check if the password is strong and matches the password policy
	// length > 8, atleast 1 upper case, atleast 1 lower case, atleast 1 symbol
	// ispasswordstrong, _ := IsPasswordStrong(tempUser.Password)
	// if ispasswordstrong == false {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong."})
	// 	return
	// }

	// Check if the user already exists.
	if DoesUserExist(tempUser.Email) {
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
		fmt.Println("err ", err.Error())
		return
	}

	SanitizedUser := models.User{
		FirstName:  tempUser.FirstName,
		LastName:   tempUser.LastName,
		Email:      tempUser.Email,
		Password:   encryptedPassword,
		UserRoleID: Role.Id, //This endpoint will be used only for customer registeration.
		CreatedAt:  time.Now(),
		IsActive:   true,
	}

	errs := models.DB.Create(&SanitizedUser).Error
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"msg": "User created successfully"})

	return
}
