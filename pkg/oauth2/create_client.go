package oauth2

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"membership-system/database"
	"membership-system/pkg/oauth2"
	"net/http"
	"strings"
)

type Client struct {
	ClientId     string   `json:"client_id"`
	ClientName   string   `form:"client_name" json:"client_name" binding:"required"`
	Scope        []string `form:"scope" json:"scope" binding:"required"`
	GrantTypes   []string `form:"grant_types" json:"grant_types" binding:"required"`
	RedirectUris []string `form:"redirect_uris" json:"redirect_uris" binding:"required"`
}

func CreateClient(c *gin.Context) {
	body := &Client{}

	if err := c.BindJSON(body); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	var client oauth2.Client
	client.Id = body.ClientId
	client.ClientName = body.ClientName
	client.Scope = strings.Join(body.Scope, " ")
	client.GrantTypes, _ = json.Marshal(body.GrantTypes)
	client.RedirectUris, _ = json.Marshal(body.RedirectUris)

	db, err := database.ConnectMysql()
	if err != nil {
		panic(err)
	}

	if result := db.Create(&client); result.Error != nil {
		c.AbortWithError(http.StatusInternalServerError, result.Error)
		return
	}

	sql, _ := db.DB()

	sql.Close()
	c.JSON(http.StatusCreated, gin.H{})
}
