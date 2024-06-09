package oauth2

import (
	"github.com/gin-gonic/gin"
	"membership-system/internal"
	"membership-system/pkg"
)

func Routes(r *gin.Engine, container internal.Container) {
	router := r.Group("/oauth2")

	// Clients
	router.GET("/clients", GetClients)
	router.GET("/client/:client_id", GetClient)
	router.POST("/client", CreateClient)
	router.PUT("/client/:client_id", UpdateClient)
	router.DELETE("/client/:client_id", DeleteClient)

	router.GET("/authorize", pkg.AuthenticationMiddleware, Authorize(container.OauthServer))
	router.POST("/token", IssueToken(container.OauthServer))
	router.POST("/login", Login)
}
