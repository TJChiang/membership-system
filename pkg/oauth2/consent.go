package oauth2

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-session/session"
	"net/http"
	"net/url"
)

func Consent(srv *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		store, err := session.Start(context.Background(), c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		var form url.Values
		if v, ok := store.Get("authorize_info"); ok {
			form = v.(url.Values)
		}
		c.Request.Form = form

		store.Delete("authorize_info")
		store.Save()

		err = srv.HandleAuthorizeRequest(c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			c.Abort()
		}
	}
}
