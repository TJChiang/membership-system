package oauth2

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"membership-system/database"
	"membership-system/internal"
	"net/http"
)

type UpdateClientRequest struct {
	ClientSecret string   `form:"client_secret" json:"client_secret"`
	ClientName   string   `form:"client_name" json:"client_name"`
	Scope        string   `form:"scope" json:"scope"`
	GrantTypes   []string `form:"grant_types" json:"grant_types"`
	RedirectUris []string `form:"redirect_uris" json:"redirect_uris"`
}

func UpdateClient(c *gin.Context) {
	id := c.Param("client_id")
	body := &UpdateClientRequest{}

	if err := c.ShouldBindBodyWithJSON(body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	var client Client

	db, err := database.ConnectMysql()
	if err != nil {
		panic(err)
	}

	sql, _ := db.DB()
	defer sql.Close()

	result := db.First(&client, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	client.ClientName = body.ClientName
	if body.ClientSecret != "" {
		client.ClientSecret, _ = internal.HashPassword(body.ClientSecret)
	}
	client.Scope = body.Scope
	client.GrantTypes, _ = json.Marshal(body.GrantTypes)
	client.RedirectUris, _ = json.Marshal(body.RedirectUris)
	client.Audience, _ = json.Marshal("[]")
	client.PostLogoutRedirectUris, _ = json.Marshal("[]")

	db.Save(&client)
	c.JSON(http.StatusOK, &client)
}
