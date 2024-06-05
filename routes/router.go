package routes

import (
	"github.com/gin-gonic/gin"
	"membership-system/pkg/oauth2"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	oauth2.Routes(router)

	return router
}
