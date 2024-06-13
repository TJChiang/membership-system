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
		c.Abort()
	})
	router.GET("/sso/static/assets", pkg.CheckGetMethod(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"path": "/sso/static/assets",
		})
	})
	router.GET("/sso/api/hello", pkg.CheckGetMethod(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"headers": c.Request.Header,
		})
	})
	router.GET("/sso/api/check-header", pkg.CheckGetMethod(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"headers": c.Request.Header,
		})
	})

	router.POST("/sso/api/do-something", pkg.CheckPostAndPutMethod(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"method":  c.Request.Method,
			"headers": c.Request.Header,
		})
	})
	router.PUT("/sso/api/do-something", pkg.CheckPostAndPutMethod(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"method":  c.Request.Method,
			"headers": c.Request.Header,
		})
	})

	router.DELETE("/sso/api/delete-something", pkg.CheckDeleteMethod(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"method":  c.Request.Method,
			"headers": c.Request.Header,
		})
	})
}
