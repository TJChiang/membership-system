package dsebd

import (
	"github.com/gin-gonic/gin"
	"membership-system/database"
	"membership-system/pkg/user"
	"net/http"
)

func MyInfo(c *gin.Context) {
	userId, _ := c.Get("user_id")
	db, err := database.ConnectMysql()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	var userModel user.User
	if result := db.First(&userModel, "id = ?", userId); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": userModel})
}
