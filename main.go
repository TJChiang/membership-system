package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"membership-system/internal"
	"membership-system/pkg"
	"membership-system/pkg/dsebd"
	"membership-system/pkg/oauth2"
	"membership-system/pkg/user"
	"net/http"
	"strings"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	container := internal.NewContainer()

	oauthServer := oauth2.Serve()

	container.OauthServer = oauthServer

	router := gin.Default()
	router.LoadHTMLGlob("./internal/templates/*.tmpl")
	router.Use(
		pkg.SetTimestamp(),
		pkg.SessionMiddleware(),
	)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/", func(c *gin.Context) {
		store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
		if err != nil {
			panic(err)
		}

		s := uuid.New().String()
		state := strings.Replace(s, "-", "", -1)

		store.Set("state", state)
		store.Save()

		u := container.OauthClient.AuthCodeURL(state)
		c.Redirect(http.StatusFound, u)
	})
	oauth2.Routes(router, container)
	dsebd.Routes(router, container)
	user.Routes(router, container)

	err := router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
