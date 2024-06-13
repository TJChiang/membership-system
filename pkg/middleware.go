package pkg

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
	redis2 "github.com/redis/go-redis/v9"
	"log"
	"membership-system/database"
	"net/http"
	"strings"
	"time"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" {
			if c.Request.URL.Path == "/dsebd/sso/resource" {
				c.Request.URL.Path = "/dsebd/sso/static/assets"
			}
			if c.Request.URL.Path == "/dsebd/me" {
				if _, err := c.Cookie("sbcookie"); err != nil {
					c.JSON(http.StatusForbidden, gin.H{"error": "sbcookie not found"})
					c.Abort()
					return
				}
			}
			if c.Request.Referer() != "www.svc.deltaww-energy.com" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Invalid referer"})
				c.Abort()
				return
			}
			if strings.HasPrefix(c.Request.URL.Path, "/dsebd/sso/api/") {
				c.Request.Header.Set("From", "hello@deltaww-energy.com")
			}
		}
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			c.Request.URL.RawQuery = ""
			if c.Request.Header.Get("X-DSEBD-AGENT") == "" {
				c.JSON(http.StatusForbidden, gin.H{"error": "X-DSEBD-AGENT header missing"})
				c.Abort()
				return
			}
			if c.Request.Header.Get("Content-Type") != "application/json" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Content-Type"})
				c.Abort()
				return
			}
		}
		if c.Request.Method == "DELETE" {
			agent := c.Request.Header.Get("X-DSEBD-AGENT")
			if agent != "AGENT_1" && agent != "AGENT_2" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Invalid X-DSEBD-AGENT value"})
				c.Abort()
				return
			}
		}
		c.Request.Header.Set("X-DSEBD-TIMESTAMP", time.Now().Format(time.RFC3339))
		c.Next()
	}
}

func CheckGetMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			if c.Request.Referer() != "www.svc.deltaww-energy.com" {
				//c.JSON(http.StatusForbidden, gin.H{"error": "Invalid referer"})
				//c.Abort()
				//return
				log.Print("Referer error, need to be http://www.svc.deltaww-energy.com ")
			}

			if strings.HasPrefix(c.Request.URL.Path, "/dsebd/sso/api/") {
				c.Request.Header.Set("From", "hello@deltaww-energy.com")
			}
		}
		c.Next()
	}
}

func CheckPostAndPutMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			if c.Request.Header.Get("X-DSEBD-AGENT") == "" {
				//c.JSON(http.StatusForbidden, gin.H{"error": "X-DSEBD-AGENT header missing"})
				//c.Abort()
				//return
				log.Println("X-DSEBD-AGENT header missing")
			}
			if c.Request.Header.Get("Content-Type") != "application/json" {
				//c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Content-Type"})
				//c.Abort()
				//return
				log.Println("Invalid Content-Type")
			}
		}
		c.Next()
	}
}

func CheckDeleteMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodDelete {
			agent := c.Request.Header.Get("X-DSEBD-AGENT")
			if agent != "alpha" && agent != "beta" {
				//c.JSON(http.StatusForbidden, gin.H{"error": "Invalid X-DSEBD-AGENT value"})
				//c.Abort()
				//return
				log.Println("Invalid X-DSEBD-AGENT value")
			}
		}
		c.Next()
	}
}

func AuthenticationMiddleware(c *gin.Context) {
	cookie, err := c.Request.Cookie("sbcookie")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sbcookie not found"})
		c.Abort()
		return
	}

	redis, err := database.ConnectRedis()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	defer redis.Close()

	log.Println("cookie value: ", cookie.Value)

	result, err := redis.Get(c.Request.Context(), cookie.Value).Result()
	if errors.Is(err, redis2.Nil) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session not found"})
		c.Abort()
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.Set("sbcookie", cookie.Value)
	c.Set("user_id", result)
	log.Println("[middleware]user_id", result)
	c.Next()
}

func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := session.Start(c.Request.Context(), c.Writer, c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to start session"})
			return
		}

		c.Next()
	}
}
