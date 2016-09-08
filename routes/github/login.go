package githubRoute

import (
	"net/http"

	"github.com/dghubble/ctxh"
	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"
	"github.com/dring1/jwt-oauth/routes"

	"golang.org/x/oauth2"
	githubOAuth2 "golang.org/x/oauth2/github"
)

type GithubLoginRoute struct {
	routes.Route
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func (ghr *GithubLoginRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	config := &oauth2.Config{
		ClientID:     ghr.ClientID,
		ClientSecret: ghr.ClientSecret,
		RedirectURL:  ghr.RedirectURL,
		Endpoint:     githubOAuth2.Endpoint,
	}

	stateConfig := gologin.DebugOnlyCookieConfig

	handler := ctxh.NewHandler(github.StateHandler(stateConfig, github.LoginHandler(config, nil)))
	handler.ServeHTTP(w, r)

}
