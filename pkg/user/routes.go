package user

import "github.com/gin-gonic/gin"

func Routes(r *gin.Engine) {
	router := r.Group("/user")

	router.POST("/register", Register)
	router.POST("/login", Login)
	router.GET("/:id")
	router.PUT("/role")
	router.GET("/protected/admin")
}
