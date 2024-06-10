package dsebd

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
	"io"
	"log"
	"membership-system/internal"
	"net/http"
	"net/url"
	"os"
	"strings"
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

		vToken, claims, err := internal.ParseAndValidateAccessToken(token.AccessToken)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		log.Println("start to get userinfo...")
		userinfo, err := getUserinfo(token.AccessToken)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":     c.Request.URL.Query(),
			"token":    token,
			"claims":   claims,
			"is_valid": vToken.Valid,
			"userinfo": userinfo,
		})
	}
}

type response struct {
	Data content `json:"data"`
}

type content struct {
	Userinfo userinfoData `json:"userinfo"`
}

type userinfoData struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  int    `json:"role"`
}

func getUserinfo(accessToken string) (*userinfoData, error) {
	appUrl := os.Getenv("APP_URL")
	userInfoEndpoint := appUrl + os.Getenv("USERINFO_ENDPOINT")

	log.Println("userinfo endpoint:", userInfoEndpoint)

	payload := strings.NewReader(url.Values{
		"client_id":     {"delta"},
		"client_secret": {"delta-secret"},
	}.Encode())

	req, err := http.NewRequest(http.MethodPost, userInfoEndpoint, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	resData := response{}
	json.Unmarshal(body, &resData)
	return &resData.Data.Userinfo, nil
}
