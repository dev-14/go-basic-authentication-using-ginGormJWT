package middleware

import (
	"fmt"
	"gingorm/controllers"
	"gingorm/models"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// the jwt middleware
var Key string

func GetAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	// var identityKey = []string{"email", "mobile"}
	// var identityKey1 = "email"
	// var identityKey2 = "mobile"
	if controllers.Flag == "email" {
		Key = "email"
	} else {
		Key = "mobile"
	}

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: Key,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				if controllers.Flag == "email" {
					return jwt.MapClaims{Key: v.Email}
				} else {
					return jwt.MapClaims{Key: v.Mobile}
				}
				//fmt.Println(v.Email)

			}

			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			if controllers.Flag == "email" {
				return &models.User{
					Email: claims[Key].(string),
				}
			} else {
				return &models.User{
					Mobile: claims[Key].(string),
				}
			}

		},
		Authenticator: controllers.Login,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*models.User); ok {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"code": code, "message": message})
		},
		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			//claims := jwt.ExtractClaims(c)
			id, _ := models.Rdb.HGet("user", "RoleID").Result()
			// claims := jwt.ExtractClaims(c)
			// fmt.Println(claims)
			// email, _ := claims["email"]
			//fmt.Println(email)
			c.JSON(code, gin.H{"code": code, "message": message, "expiry": time, "AccessLevel": id})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			email := c.GetString("user_email")
			models.Rdb.Del(email)
			fmt.Println("Redis Cleared")
			c.JSON(code, gin.H{
				"message": "logged out successfully",
			})
		},
		RefreshResponse: func(*gin.Context, int, string, time.Time) {
		},

		TokenLookup:   "header: Authorization, query: token, cookie: token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		// HTTPStatusMessageFunc: func(e error, c *gin.Context) string {
		// },
		PrivKeyFile:       "",
		PubKeyFile:        "",
		SendCookie:        true,
		SecureCookie:      true,
		CookieHTTPOnly:    true,
		SendAuthorization: true,
		DisabledAbort:     false,
		CookieName:        "token",
	})
	if err != nil {
		return nil, err
	}
	return authMiddleware, nil
}
