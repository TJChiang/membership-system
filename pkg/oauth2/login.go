package oauth2

import (
	"github.com/gin-gonic/gin"
	"membership-system/database"
	"membership-system/internal"
	model "membership-system/pkg/user"
	"net/http"
	"strings"
	"time"
)

type loginRequest struct {
	Username string `form:"username" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var credentials loginRequest
	if err := c.ShouldBind(&credentials); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	var user model.User
	db, err := database.ConnectMysql()
	if err != nil {
		panic(err)
	}

	sql, _ := db.DB()
	defer sql.Close()

	username := strings.ToLower(strings.TrimSpace(credentials.Username))
	if result := db.First(&user, "email = ?", username); result.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	if !internal.CheckPasswordHash(credentials.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	sessionKey, err := internal.GenerateRandomStringURLSafe(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rdb, err := database.ConnectRedis()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to redis"})
		return
	}

	ttl := 30 * time.Minute
	expiration := time.Now().Add(ttl)
	err = rdb.Set(c, sessionKey, user.Id, ttl).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error storing session"})
		return
	}

	c.SetCookie(
		"sbcookie",
		sessionKey,
		int(expiration.Unix()),
		"/",
		".local",
		false,
		true,
	)

	// redirect to Authorization page
	c.Redirect(http.StatusFound, "/dsebd/consent")
}
