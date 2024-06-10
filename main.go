package main

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	goauth2 "golang.org/x/oauth2"
	"membership-system/internal"
	"membership-system/pkg/dsebd"
	"membership-system/pkg/oauth2"
	"net/http"
)

var (
	goauth2Config = goauth2.Config{
		ClientID:     "delta",
		ClientSecret: "delta-secret",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:8080/dsebd/callback",
		Endpoint: goauth2.Endpoint{
			AuthURL:  "http://localhost:8080/oauth2/authorize",
			TokenURL: "http://localhost:8080/oauth2/token",
		},
	}
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	var container internal.Container

	oauthServer := oauth2.Serve()

	container.OauthServer = oauthServer

	router := gin.Default()
	router.LoadHTMLGlob("./internal/templates/*.tmpl")
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/", func(c *gin.Context) {
		safe, err := internal.GenerateRandomStringURLSafe(16)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		s256 := sha256.Sum256([]byte("s256example"))
		challenge := base64.URLEncoding.EncodeToString(s256[:])

		u := goauth2Config.AuthCodeURL(safe,
			goauth2.SetAuthURLParam("code_challenge", challenge),
			goauth2.SetAuthURLParam("code_challenge_method", "S256"))
		c.Redirect(http.StatusFound, u)
	})
	oauth2.Routes(router, container)
	dsebd.Routes(router)

	err := router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
