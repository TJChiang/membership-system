package oauth2

import (
	"log"
	"net/http"
	"os"

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
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uid, ok := store.Get("oauth2_subject")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}
		store.Set("authorize_info", r.Form)
		store.Save()

		http.Redirect(w, r, "/dsebd/login", http.StatusFound)
		return
	}

	log.Println("subject has logged in: ", uid)
	http.Redirect(w, r, "/dsebd/me", http.StatusFound)
	return
}
