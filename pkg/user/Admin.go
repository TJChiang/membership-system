package user

import (
	"github.com/gin-gonic/gin"
	"membership-system/database"
	"net/http"
)

func AdminInfo(c *gin.Context) {
	uid, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "User info not found in context",
		})
		return
	}

	var user User

	db, err := database.ConnectMysql()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	sql, _ := db.DB()
	defer sql.Close()

	if result := db.First(&user, "id = ?", uid); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	if user.Role != Admin.Value() {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "User info not found in context",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
