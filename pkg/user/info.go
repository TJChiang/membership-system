package user

import (
	"github.com/gin-gonic/gin"
	"membership-system/database"
	"net/http"
)

func Info(c *gin.Context) {
	uid := c.Param("id")

	suid, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}

	var suser, user User

	db, err := database.ConnectMysql()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	sql, _ := db.DB()
	defer sql.Close()

	if result := db.First(&suser, "id = ?", suid); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if uid == suid {
		c.JSON(http.StatusOK, gin.H{
			"user": suser,
		})
		return
	}

	if suser.Role != Admin.Value() {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	if result := db.First(&user, "id = ?", uid); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
