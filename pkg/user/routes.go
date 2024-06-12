package user

import (
	"github.com/gin-gonic/gin"
	"membership-system/internal"
	"membership-system/pkg"
)

func Routes(r *gin.Engine, container *internal.Container) {
	router := r.Group("/user")

	router.GET("/:id", pkg.AuthenticationMiddleware, Info)
	router.PUT("/role", pkg.AuthenticationMiddleware)
	router.GET("/protected/admin", pkg.AuthenticationMiddleware, AdminInfo)
}
