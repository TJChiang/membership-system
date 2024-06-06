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
		return
	}

	sql, _ := db.DB()
	defer sql.Close()

	result := db.Find(&clients)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": &clients,
	})
}
