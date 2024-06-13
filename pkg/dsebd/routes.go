package dsebd

import (
	"github.com/gin-gonic/gin"
	"membership-system/internal"
	"membership-system/pkg"
	"net/http"
)

func Routes(r *gin.Engine, container *internal.Container) {
	router := r.Group("/dsebd")

	router.GET("/register", RegisterPage)
	router.POST("/register", Register)
	router.GET("/login", LoginPage)
	router.GET("/consent", ConsentPage)
	router.GET("/me", pkg.AuthenticationMiddleware, MyInfo)
	router.GET("/callback", OAuthCallback(container))

	router.GET("/sso/resource", func(c *gin.Context) {
		c.Request.URL.Path = "/dsebd/sso/static/assets"
		r.HandleContext(c)
	})
	router.GET("/sso/static/assets", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"path": "/sso/static/assets",
		})
	})

	router.GET("/sso/api/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"headers": c.Request.Header,
		})
	})
	router.GET("/sso/api/check-header", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"headers": c.Request.Header,
		})
	})
}
