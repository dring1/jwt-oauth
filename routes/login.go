package routes

import (
	"log"
	"os"

	"github.com/dghubble/ctxh"
	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"
	"github.com/dghubble/sessions"
	"github.com/dring1/orm/controllers"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	githubOAuth2 "golang.org/x/oauth2/github"
)

type Config struct {
	GithubClientID     string
	GithubClientSecret string
}

const (
	sessionName    = "example-github-app"
	sessionSecret  = "example cookie signing secret"
	sessionUserKey = "githubID"
)

var c *Config

var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

func init() {
	if val := os.Getenv("GITHUB_CLIENT_ID"); val == "" {
		log.Fatal("GITHUB_CLIENT_ID NOT SET")
	}
	if val := os.Getenv("GITHUB_CLIENT_SECRET"); val == "" {
		log.Fatal("GITHUB_CLIENT_SECRET NOT SET")
	}

	c = &Config{
		GithubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GithubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	}
}

func LoginRoute(r *mux.Router) *mux.Router {
	oauth2Config := &oauth2.Config{
		ClientID:     c.GithubClientID,
		ClientSecret: c.GithubClientSecret,
		RedirectURL:  "http://localhost:8080/github/callback",
		Endpoint:     githubOAuth2.Endpoint,
	}
	// state param cookies require HTTPS by default; disable for localhost development
	stateConfig := gologin.DebugOnlyCookieConfig
	r.Handle("/github/login", ctxh.NewHandler(github.StateHandler(stateConfig, github.LoginHandler(oauth2Config, nil))))
	r.Handle("/github/callback", ctxh.NewHandler(github.StateHandler(stateConfig, github.CallbackHandler(oauth2Config, controllers.Login(), nil))))
	return r
}
