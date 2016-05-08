package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dghubble/ctxh"
	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"
	"github.com/dghubble/sessions"
	"github.com/dring1/orm/controllers"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
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
	r.Handle("/github/callback", ctxh.NewHandler(github.StateHandler(stateConfig, github.CallbackHandler(oauth2Config, issueSession(), nil))))
	return r
}

func issueSession() ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		controllers.Login(w, req)
		// 2. Implement a success handler to issue some form of session
		session := sessionStore.New(sessionName)
		session.Values[sessionUserKey] = *githubUser.ID
		session.Save(w)
		fmt.Println(session.Values)
		http.Redirect(w, req, "/profile", http.StatusFound)
	}
	return ctxh.ContextHandlerFunc(fn)
}
