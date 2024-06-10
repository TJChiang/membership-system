package dsebd

import (
	"github.com/gin-gonic/gin"
	"membership-system/pkg"
)

func Routes(r *gin.Engine) {
	router := r.Group("/dsebd")

	router.GET("/register", RegisterPage)
	router.POST("/register", Register)
	router.GET("/login", LoginPage)
	router.GET("/consent", ConsentPage)
	router.GET("/me", pkg.AuthenticationMiddleware, MyInfo)
	router.GET("/callback", OAuthCallback)
}
