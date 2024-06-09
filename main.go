package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"membership-system/internal"
	"membership-system/pkg/dsebd"
	"membership-system/pkg/oauth2"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	var container internal.Container

	oauthServer := oauth2.Serve()

	container.OauthServer = oauthServer

	router := gin.Default()
	router.LoadHTMLGlob("/**/*/templates/*.tmpl")
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	oauth2.Routes(router, container)
	dsebd.Routes(router)

	err := router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
