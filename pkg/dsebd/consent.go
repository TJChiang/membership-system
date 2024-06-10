package dsebd

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
)

func ConsentPage(c *gin.Context) {
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	_, ok := store.Get("oauth2_subject")
	if !ok {
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "consent.tmpl", gin.H{
		"title":       "Consent Page",
		"consent_url": "/oauth2/consent",
	})
}
