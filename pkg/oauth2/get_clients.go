package oauth2

import (
	"github.com/gin-gonic/gin"
	"membership-system/database"
	"net/http"
)

func GetClients(c *gin.Context) {
	var clients []Client

	db, err := database.ConnectMysql()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	result := db.First(&clients)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": &clients,
	})
}
