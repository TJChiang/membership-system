package dsebd

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ConsentPage(c *gin.Context) {
	c.HTML(http.StatusOK, "consent.tmpl", gin.H{})
}
