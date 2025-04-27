package api

import (
	"github.com/gin-gonic/gin"
	"personal-secretary-user-ap/internal/api/handler"
	"personal-secretary-user-ap/internal/api/handler/register"
)

func GetRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", handler.Home)

	v1 := router.Group("/v1")
	v1.POST("/register", register.Register)

	/*v1.Use(middleware.AuthMiddleware())
	{
		v1.GET("/init", handler.Init)
	}*/

	return router
}
