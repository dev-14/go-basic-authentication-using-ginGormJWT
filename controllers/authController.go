package controllers

import (
	"gingorm/models"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	//"github.com/gofiber/fiber/v2"

	//"github.com/gofiber/fiber"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

type RegisterNew struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterNew
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: password,
	}

	models.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func Login(c *gin.Context) {
	var data map[string]string

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
	}

	var user models.User

	models.DB.Where("email=?", data["email"]).First(&user)
	if user.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
	} else if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "wrong password",
		})
	} else {
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    strconv.Itoa(int(user.Id)),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		})

		token, err := claims.SignedString([]byte(SecretKey))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "could not login",
			})
		}
		// _ = token

		// c.SetCookie("jwt", token, 60*60*24, "", "", true, true)
		http.SetCookie(c.Writer, &http.Cookie{
			Name:    user.Name,
			Value:   token,
			Expires: time.Now().Add(time.Hour * 24),
		})

		c.JSON(http.StatusOK, gin.H{
			"message": "successfully logged in",
			"token":   token,
		})
	}

}

func User(c *gin.Context) {
	var user models.User
	// user := models.User
	cookie, err := c.Request.Cookie(user.Name)
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
