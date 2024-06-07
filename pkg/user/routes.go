package user

import (
	"github.com/gin-gonic/gin"
	"membership-system/pkg"
)

func Routes(r *gin.Engine) {
	router := r.Group("/user")

	router.GET("/:id", pkg.AuthenticationMiddleware)
	router.PUT("/role", pkg.AuthenticationMiddleware)
	router.GET("/protected/admin", pkg.AuthenticationMiddleware)
}
