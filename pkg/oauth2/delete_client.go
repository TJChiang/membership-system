package oauth2

import (
	"github.com/gin-gonic/gin"
	"log"
	"membership-system/database"
	"net/http"
)

func DeleteClient(c *gin.Context) {
	id := c.Param("client_id")

	var client Client

	db, err := database.ConnectMysql()
	if err != nil {
		panic(err)
	}

	sql, _ := db.DB()
	defer sql.Close()

	if result := db.First(&client, "id = ?", id); result.Error != nil {
		log.Println("Failed to delete the client: ", id, " for", result.Error.Error())
		c.Status(http.StatusOK)
		return
	}

	db.Delete(&client)
	c.Status(http.StatusOK)
}
