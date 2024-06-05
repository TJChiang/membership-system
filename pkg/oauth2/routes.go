package oauth2

import (
	"github.com/gin-gonic/gin"
	"membership-system/api/oauth2"
)

func Routes(r *gin.Engine) {
	router := r.Group("/oauth2")
	// Clients
	router.GET("/clients")
	router.GET("/client/:client_id")
	router.POST("/client", oauth2.CreateClient)
	router.PUT("/client/:client_id")
	router.DELETE("/client/:client_id")

	router.GET("/auth")
}
