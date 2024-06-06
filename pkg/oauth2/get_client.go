package oauth2

import (
	"github.com/gin-gonic/gin"
	"membership-system/database"
	"net/http"
)

func GetClient(c *gin.Context) {
	id := c.Param("client_id")

	var client Client

	db, err := database.ConnectMysql()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	sql, _ := db.DB()
	defer sql.Close()

	result := db.First(&client, "id = ?", id)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": &client,
	})
}
