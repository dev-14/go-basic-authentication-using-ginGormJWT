package routes

import (
	"gingorm/controllers"
	//_ "gingorm/docs"
	middlewares "gingorm/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func PublicEndpoints(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	// Generate public endpoints - [ signup] - api/v1/signup

	r.POST("/register", controllers.Register)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/login", authMiddleware.LoginHandler)
}

func GetRouter(router chan *gin.Engine) {
	gin.ForceConsoleColor()
	r := gin.Default()

	r.Use(cors.Default())

	authMiddleware, _ := middlewares.GetAuthMiddleware()

	// Create a BASE_URL - /api/v1
	v1 := r.Group("/api/v1/")
	PublicEndpoints(v1, authMiddleware)
	//AuthenticatedEndpoints(v1.Group("auth"), authMiddleware)
	router <- r
}
