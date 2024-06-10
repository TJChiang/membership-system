package examples

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	oredis "github.com/go-oauth2/redis/v4"
	"github.com/go-redis/redis/v8"
	"github.com/go-session/session"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	goauth2 "golang.org/x/oauth2"
	"log"
	"membership-system/internal"
	"membership-system/pkg"
	"membership-system/pkg/dsebd"
	"net/http"
	"net/url"
	"os"
)

var (
	testConfig = goauth2.Config{
		ClientID:     "delta",
		ClientSecret: "delta-secret",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:8080/dsebd/callback",
		Endpoint: goauth2.Endpoint{
			AuthURL:  "http://localhost:8080/auth",
			TokenURL: "http://localhost:8080/token",
		},
	}
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	manager := manage.NewDefaultManager()
	manager.MapTokenStorage(oredis.NewRedisStore(&redis.Options{
		Addr: os.Getenv("REDIS_HOST"),
		DB:   15,
	}))
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("membership-secret"), jwt.SigningMethodHS512))

	clientStore := store.NewClientStore()
	clientStore.Set("delta", &models.Client{
		ID:     "delta",
		Secret: "delta-secret",
		Domain: "http://localhost:8080",
	})
	manager.MapClientStorage(clientStore)

	srv := server.NewServer(server.NewConfig(), manager)

	srv.SetUserAuthorizationHandler(userAuthorizeHandler)
	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})
	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	router := gin.Default()

	router.LoadHTMLGlob("./internal/templates/*.tmpl")
	router.Use(pkg.SessionMiddleware())
	router.GET("/login", loginHandler)
	router.POST("/login", loginPostHandler)
	router.GET("/auth", func(c *gin.Context) {
		store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var form url.Values
		if v, ok := store.Get("ReturnUri"); ok {
			form = v.(url.Values)
		}
		c.Request.Form = form

		store.Delete("ReturnUri")
		store.Save()

		err = srv.HandleAuthorizeRequest(c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})
	router.POST("/token", func(c *gin.Context) {
		err := srv.HandleTokenRequest(c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	router.GET("/", func(c *gin.Context) {
		safe, err := internal.GenerateRandomStringURLSafe(16)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		s256 := sha256.Sum256([]byte("s256example"))
		challenge := base64.URLEncoding.EncodeToString(s256[:])

		u := testConfig.AuthCodeURL(safe,
			goauth2.SetAuthURLParam("code_challenge", challenge),
			goauth2.SetAuthURLParam("code_challenge_method", "S256"))
		c.Redirect(http.StatusFound, u)
	})

	dsebd.Routes(router)
	router.Run(":8080")
}

func loginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"title":     "Login",
		"login_url": "/login",
	})
}

func loginPostHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "user" && password == "pass" {
		store, _ := session.Start(c.Request.Context(), c.Writer, c.Request)
		store.Set("LoggedInUserID", "user_id")
		store.Save()
		c.Redirect(http.StatusFound, "/auth")
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return "", err
	}

	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}
		store.Set("ReturnUri", r.Form)
		store.Save()
		http.Redirect(w, r, "/login", http.StatusFound)
		return "", nil
	}

	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	return
}
