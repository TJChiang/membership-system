package oauth2

import (
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/server"
	"membership-system/database"
	model "membership-system/pkg/user"
	"net/http"
	"time"
)

func Userinfo(srv *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := srv.ValidationBearerToken(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		var user model.User
		uid := token.GetUserID()
		if err := getUserInfo(&user, uid); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		userInfo := map[string]interface{}{
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		}

		data := map[string]interface{}{
			"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
			"client_id":  token.GetClientID(),
			"userinfo":   userInfo,
		}

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}

func getUserInfo(user *model.User, uid string) error {
	db, err := database.ConnectMysql()
	if err != nil {
		return err
	}

	sql, err := db.DB()
	if err != nil {
		return err
	}

	defer sql.Close()

	if result := db.First(&user, "id = ?", uid); result.Error != nil {
		return result.Error
	}

	return nil
}
