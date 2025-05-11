package api

import (
	"github.com/gin-gonic/gin"
	"personal-secretary-user-ap/internal/api/handler"
	"personal-secretary-user-ap/internal/api/handler/login"
	"personal-secretary-user-ap/internal/api/handler/register"
	"personal-secretary-user-ap/internal/api/middleware"
	"personal-secretary-user-ap/internal/application"
)

func GetRouter() *gin.Engine {
	// Get application instance
	app, err := application.GetInstance(nil)
	if err != nil {
		panic("Failed to get application instance: " + err.Error())
	}

	// Get JWT configuration
	jwtSecret := app.GetConfig().JWT.Secret
	router := gin.Default()

	router.Use(middleware.ResponseLogger())

	router.GET("/", handler.Home)

	v1 := router.Group("/v1")
	v1.POST("/login", login.Login)
	v1.POST("/register", register.Register)

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// token management

		// Add other protected routes here
		// Example: protected.GET("/user/profile", userHandler.GetProfile)
	}

	return router
}
