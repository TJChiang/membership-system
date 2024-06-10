package oauth2

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-session/session"
)

func Consent(srv *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var form url.Values
		v, ok := store.Get("authorize_info")
		if !ok {
			c.Redirect(http.StatusFound, "/")
			return
		}
		form = v.(url.Values)
		c.Request.Form = form

		store.Delete("authorize_info")
		store.Save()

		err = srv.HandleAuthorizeRequest(c.Writer, c.Request)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
	}
}
