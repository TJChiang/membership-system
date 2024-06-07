package dsebd

import (
	"github.com/gin-gonic/gin"
	"membership-system/database"
	"membership-system/internal"
	user2 "membership-system/pkg/user"
	"net/http"
	"strings"
)

type RegisterRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	body := &RegisterRequest{}
	if err := c.ShouldBindBodyWithJSON(body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user user2.User
	hashPassword, err := internal.HashPassword(strings.TrimSpace(body.Password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.Name = strings.TrimSpace(body.Username)
	user.Email = strings.ToLower(strings.TrimSpace(body.Email))
	user.Password = hashPassword
	user.Role = user2.Member.Value()

	db, err := database.ConnectMysql()
	if err != nil {
		panic(err)
	}

	sql, _ := db.DB()
	defer sql.Close()

	if result := db.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.Status(http.StatusCreated)
}
