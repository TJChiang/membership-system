package oauth2

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"membership-system/database"
	"membership-system/internal"
	"net/http"
)

type CreateClientRequest struct {
	ClientId     string   `form:"client_id" json:"client_id" binding:"required"`
	ClientSecret string   `form:"client_secret" json:"client_secret" binding:"required"`
	ClientName   string   `form:"client_name" json:"client_name" binding:"required"`
	Scope        string   `form:"scope" json:"scope" binding:"required"`
	GrantTypes   []string `form:"grant_types" json:"grant_types" binding:"required"`
	RedirectUris []string `form:"redirect_uris" json:"redirect_uris" binding:"required"`
}

func CreateClient(c *gin.Context) {
	body := &CreateClientRequest{}

	if err := c.ShouldBindBodyWithJSON(body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	var client Client
	client.Id = body.ClientId
	client.ClientName = body.ClientName
	client.ClientSecret, _ = internal.HashPassword(body.ClientSecret)
	client.Scope = body.Scope
	client.GrantTypes, _ = json.Marshal(body.GrantTypes)
	client.RedirectUris, _ = json.Marshal(body.RedirectUris)
	client.Audience, _ = json.Marshal("[]")
	client.PostLogoutRedirectUris, _ = json.Marshal("[]")

	db, err := database.ConnectMysql()
	if err != nil {
		panic(err)
	}

	if result := db.Create(&client); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	sql, _ := db.DB()

	sql.Close()
	c.JSON(http.StatusCreated, gin.H{})
}
