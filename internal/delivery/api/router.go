package api

import (
	"github.com/Myakun/personal-secretary-user-api/internal/delivery/api/handler"
	registerHandler "github.com/Myakun/personal-secretary-user-api/internal/delivery/api/handler/register"
	"github.com/Myakun/personal-secretary-user-api/internal/delivery/api/middleware"
	"github.com/gin-gonic/gin"
)

type RouterHandlers struct {
	RegisterHandler *registerHandler.RegisterHandler
}

func GetRouter(handlers *RouterHandlers) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.ResponseLogger())

	router.GET("/", handler.Home)

	v1 := router.Group("/v1")
	v1.POST("/register", handlers.RegisterHandler.Register)

	// Protected routes
	//	protected := v1.Group("")
	//	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// token management

		// Add other protected routes here
		// Example: protected.GET("/user/profile", userHandler.GetProfile)
	}

	//router := gin.Default()
	//router.GET("/", handler.Home)

	return router
}
