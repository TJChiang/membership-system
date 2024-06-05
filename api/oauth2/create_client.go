package oauth2

import "github.com/gin-gonic/gin"

type Client struct {
	ClientName string `form:"client_name" json:"client_name" binding:"required"`
	Scope      string `form:"scope" json:"scope" binding:"required"`
}

func CreateClient(c *gin.Context) {}
