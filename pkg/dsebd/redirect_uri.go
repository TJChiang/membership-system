package dsebd

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
	"log"
	"membership-system/internal"
	"net/http"
)

func OAuthCallback(con *internal.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		errorMsg, ee := c.GetQuery("error")
		if ee {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errorMsg,
				"data":  c.Request.URL.Query(),
			})
			return
		}

		store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		sstate, exists := store.Get("state")
		if !exists {
			log.Println("state not found")
			c.Redirect(http.StatusFound, "/")
			return
		}

		state, ok := c.GetQuery("state")
		if !ok {
			c.AbortWithError(http.StatusBadRequest, errors.New("missing state"))
			return
		}

		log.Println("state:", state)
		log.Println("session state:", sstate)
		if sstate != state {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid state"))
			return
		}

		code := c.Query("code")
		if code == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("missing code"))
			return
		}

		log.Println("code: ", code)
		token, err := con.OauthClient.Exchange(c.Request.Context(), code)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  c.Request.URL.Query(),
			"token": token,
		})
	}
}
