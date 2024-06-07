package dsebd

import (
	"github.com/gin-gonic/gin"
	"membership-system/pkg"
)

func Routes(r *gin.Engine) {
	router := r.Group("/dsebd")

	router.POST("/register", Register)
	router.POST("/login", Login)
	router.GET("/me", pkg.AuthenticationMiddleware, MyInfo)
}
