package oauth2

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-session/session"
	"net/http"
)

type request struct {
	ClientId     string `query:"client_id"  binding:"required"`
	ResponseType string `query:"response_type"  binding:"required"`
	State        string `query:"state"  binding:"required"`
	Scope        string `query:"scope"  binding:"required"`
	RedirectUri  string `query:"redirect_uri"  binding:"required"`
}

func Authorize(srv *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		store, err := session.Start(context.Background(), c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		store.Delete("authorize_info")
		store.Save()

		if c.Request.Form == nil {
			c.Request.ParseForm()
		}

		err = srv.HandleAuthorizeRequest(c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			c.Abort()
			return
		}
	}
}
