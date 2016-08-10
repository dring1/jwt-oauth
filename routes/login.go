package routes

import (
	"log"
	"net/http"

	"github.com/dghubble/ctxh"
	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"
	"github.com/dghubble/sessions"
	"github.com/dring1/jwt-oauth/controllers"
	"github.com/dring1/jwt-oauth/models"
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

type LoginRoute struct {
	GitHubClientID     string
	GitHubClientSecret string
	Oauth2Config       *oauth2.Config
	RedirectURL        string
	LoginHandler       http.Handler
	CallbackHandler    http.Handler
}

func (r *LoginRoute) GenHttpHandlers() ([]*R, error) {
	if r.Oauth2Config == nil {
		r.Oauth2Config = &oauth2.Config{
			ClientID:     r.GitHubClientID,
			ClientSecret: r.GitHubClientSecret,
			RedirectURL:  r.RedirectURL,
			Endpoint:     githubOAuth2.Endpoint,
		}
	}
	stateConfig := gologin.DebugOnlyCookieConfig

	if r.LoginHandler == nil {
		r.LoginHandler = ctxh.NewHandler(github.StateHandler(stateConfig, github.LoginHandler(r.Oauth2Config, nil)))
	}
	if r.CallbackHandler == nil {
		r.CallbackHandler = ctxh.NewHandler(github.StateHandler(stateConfig, github.CallbackHandler(r.Oauth2Config, controllers.Login(func(u *models.User) {}), nil)))
	}

	routes := []*R{
		&R{
			Path:    "/github/login",
			Methods: []string{"GET", "POST"},
			Handler: r.LoginHandler,
		},
		&R{
			Path:    "/github/callback",
			Methods: []string{"GET", "POST"},
			Handler: r.CallbackHandler,
		},
	}
	return routes, nil
}

func OldHttpHandler(r *mux.Router, loginHandler http.Handler, callbackHandler http.Handler, oauth2Config *oauth2.Config) *mux.Router {
	log.Println("LOGIN ROUTE")
	if oauth2Config == nil {
		log.Println("LOGIN config")
		oauth2Config = &oauth2.Config{
			ClientID:     c.GithubClientID,
			ClientSecret: c.GithubClientSecret,
			RedirectURL:  "http://localhost:8080/github/callback",
			Endpoint:     githubOAuth2.Endpoint,
		}
	}
	stateConfig := gologin.DebugOnlyCookieConfig
	// log.Println(loginHandler, callbackHandler)
	// state param cookies require HTTPS by default; disable for localhost development
	if loginHandler == nil {
		loginHandler = ctxh.NewHandler(github.StateHandler(stateConfig, github.LoginHandler(oauth2Config, nil)))
	}
	if callbackHandler == nil {
		callbackHandler = ctxh.NewHandler(github.StateHandler(stateConfig, github.CallbackHandler(oauth2Config, controllers.Login(func(u *models.User) {}), nil)))
	}
	r.Handle("/github/login", loginHandler)
	r.Handle("/github/callback", callbackHandler)
	return r
}
