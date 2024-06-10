package dsebd

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func OAuthCallback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": c.Request.URL.Query(),
	})
}
