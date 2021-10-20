package controllers

import (
	"gingorm/models"
	//"errors"
	"fmt"
	"log"

	//"unicode"

	//jwt "github.com/appleboy/gin-jwt/v2"
	//"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CheckCredentials(useremail, userpassword string, db *gorm.DB) bool {
	// db := c.MustGet("db").(*gorm.DB)
	// var db *gorm.DB
	var User models.User
	// Store user supplied password in mem map
	var expectedpassword string
	// check if the email exists
	err := db.Where("email = ?", useremail).First(&User).Error
	if err == nil {
		// User Exists...Now compare his password with our password
		expectedpassword = User.Password
		if err = bcrypt.CompareHashAndPassword([]byte(expectedpassword), []byte(userpassword)); err != nil {
			// If the two passwords don't match, return a 401 status
			log.Println("User is Not Authorized")
			return false
		}
		// User is AUthenticates, Now set the JWT Token
		fmt.Println("User Verified")
		return true
	} else {
		// returns an empty array, so simply pass as not found, 403 unauth
		log.Fatal("ERR ", err)

	}
	return false
}
