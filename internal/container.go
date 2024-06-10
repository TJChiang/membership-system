package internal

import (
	"github.com/go-oauth2/oauth2/v4/server"
	"golang.org/x/oauth2"
	"os"
)

type Container struct {
	OauthServer *server.Server
	OauthClient oauth2.Config
}

func NewContainer() *Container {
	appUrl := os.Getenv("APP_URL")
	return &Container{
		OauthClient: oauth2.Config{
			ClientID:     "delta",
			ClientSecret: "delta-secret",
			Scopes:       []string{"all"},
			RedirectURL:  appUrl + "/dsebd/callback",
			Endpoint: oauth2.Endpoint{
				AuthURL:  appUrl + "/oauth2/authorize",
				TokenURL: appUrl + "/oauth2/token",
			},
		},
	}
}
