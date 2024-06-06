package oauth2

import (
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	router := r.Group("/oauth2")
	// Clients
	router.GET("/clients", GetClients)
	router.GET("/client/:client_id", GetClient)
	router.POST("/client", CreateClient)
	router.PUT("/client/:client_id")
	router.DELETE("/client/:client_id")

	router.GET("/auth")
}
