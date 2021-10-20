package controllers

import (
	"gingorm/models"
	"net/http"

	// "github.com/dgrijalva/jwt-go"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	//"github.com/gofiber/fiber/v2"

	//"github.com/gofiber/fiber"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

// Register godoc
// @Summary Creates a new User in the database
// @Description API Endpoint to register the user.
// @Router /register [post]
// @Tags user
// @Accept json
// @Produce json
type RegisterNew struct {
	FirstName string `json:"firstname" binding:"required"`
	LastName  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterNew
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	user := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  string(password),
	}

	models.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

type login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Login godoc
// @Summary Lets the user Login to the system
// @Description API Endpoint to let the user login.
// @Router /login [post]
// @Tags user
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

// User godoc
// @Summary Returns the logged in user
// @Description API Endpoint to return the name of the user that has logged in to the system.
// @Router /user [get]
// @Tags user
// @Accept json
// @Produce json
func User(c *gin.Context) {
	var user models.User
	// user := models.User
	cookie, err := c.Request.Cookie(user.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "cookie not found",
			"error":   err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"cookie": cookie.Value,
			"user":   cookie.Name,
		})
	}

	//claims := cookie.Claims.(*jwt.StandardClaims)

	//var user models.User

	// models.DB.Where("id=?", cookie.Issuer).First(&user)

	// c.JSON(http.StatusOK, gin.H{
	// 	"data": user,
	// })
}

// func Logout(c *gin.Context) {
// 	var user models.User

// 	// http.SetCookie(c.Writer, &http.Cookie{
// 	// 	Name:   user.Name,
// 	// 	MaxAge: -1,
// 	// })

// 	cookie, err := c.Request.Cookie(user.Name)

// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"message": "User not logged in",
// 		})
// 	}
// 	c.JSON(http.StatusFound, gin.H{
// 		"message": "cookie found",
// 		"user":    cookie.Name,
// 		"maxage":  cookie.MaxAge,
// 	})
// 	//cookie.Expires = time.Now().Add(-time.Hour)
// 	cookie.MaxAge = -1
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Successfully logged out",
// 		"maxage":  cookie.MaxAge,
// 	})
// }
