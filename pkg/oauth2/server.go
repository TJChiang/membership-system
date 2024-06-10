package oauth2

import (
	"context"
	"errors"
	oerrors "github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	oredis "github.com/go-oauth2/redis/v4"
	"github.com/go-redis/redis/v8"
	"github.com/go-session/session"
	"github.com/golang-jwt/jwt"
	redis2 "github.com/redis/go-redis/v9"
	"log"
	"membership-system/database"
	"net/http"
	"os"
)

func Serve() *server.Server {
	manager := manage.NewDefaultManager()

	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	manager.MapTokenStorage(oredis.NewRedisStore(&redis.Options{
		Addr: os.Getenv("REDIS_HOST"),
		DB:   15,
	}))
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("membership-secret"), jwt.SigningMethodHS512))
	//manager.MapAccessGenerate(generates.NewAccessGenerate())

	clientStore := store.NewClientStore()
	err := clientStore.Set("delta", &models.Client{
		ID:     "delta",
		Secret: "delta-secret",
		Domain: "http://localhost:8080",
	})
	if err != nil {
		panic(err)
	}
	err = clientStore.Set("alpha", &models.Client{
		ID:     "alpha",
		Secret: "alpha-secret",
		Domain: "http://alpha.local",
	})
	if err != nil {
		panic(err)
	}
	err = clientStore.Set("beta", &models.Client{
		ID:     "beta",
		Secret: "beta-secret",
		Domain: "http://beta.local",
	})
	if err != nil {
		panic(err)
	}
	manager.MapClientStorage(clientStore)

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	srv.SetInternalErrorHandler(func(err error) (re *oerrors.Response) {
		log.Println("OAuth2 Server In Internal Error:", err.Error())
		return
	})
	srv.SetResponseErrorHandler(func(re *oerrors.Response) {
		log.Println("OAuth2 Server In Response Error:", re.Error.Error())
	})

	return srv
}

// auth endpoint 之後，確認登入身份
func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie("sbcookie")
	if err != nil {
		if r.Form == nil {
			r.ParseForm()
		}
		store.Set("authorize_info", r.Form)
		store.Save()

		w.Header().Set("Location", "/dsebd/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	redis, err := database.ConnectRedis()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer redis.Close()

	_, err = redis.Get(r.Context(), cookie.Value).Result()
	if errors.Is(err, redis2.Nil) {
		if r.Form == nil {
			r.ParseForm()
		}
		store.Set("authorize_info", r.Form)
		store.Save()

		w.Header().Set("Location", "/dsebd/login")
		w.WriteHeader(http.StatusFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/dsebd/me")
	w.WriteHeader(http.StatusFound)
	return
}
