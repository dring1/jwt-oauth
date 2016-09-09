package routes

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/context"

	"github.com/dghubble/ctxh"
	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"
	"github.com/dghubble/sessions"

	"golang.org/x/oauth2"
	githubOAuth2 "golang.org/x/oauth2/github"
)

const (
	sessionName    = "example-github-app"
	sessionSecret  = "example cookie signing secret"
	sessionUserKey = "githubID"
)

var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

type GithubLoginRoute struct {
	Route
	ClientID     string
	ClientSecret string
	RedirectURL  string
	UserService  users.Service
}

func (ghr *GithubLoginRoute) NewHandler() (*R, error) {
	config := &oauth2.Config{
		ClientID:     ghr.ClientID,
		ClientSecret: ghr.ClientSecret,
		RedirectURL:  ghr.RedirectURL,
		Endpoint:     githubOAuth2.Endpoint,
	}

	stateConfig := gologin.DebugOnlyCookieConfig

	handler := ctxh.NewHandler(github.StateHandler(stateConfig, github.LoginHandler(config, nil)))
	return &R{
		Path:    ghr.Path,
		Methods: ghr.Methods,
		Handler: handler,
	}, nil

}

type GithubCallbackRoute struct {
	Route
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func (ghr *GithubCallbackRoute) NewHandler() (*R, error) {
	config := &oauth2.Config{
		ClientID:     ghr.ClientID,
		ClientSecret: ghr.ClientSecret,
		RedirectURL:  ghr.RedirectURL,
		Endpoint:     githubOAuth2.Endpoint,
	}
	stateConfig := gologin.DebugOnlyCookieConfig
	handler := ctxh.NewHandler(github.StateHandler(stateConfig, github.CallbackHandler(config, issueSession(), nil)))

	return &R{
		Path:    ghr.Path,
		Methods: ghr.Methods,
		Handler: handler,
	}, nil
}

func issueSession() ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 2. Implement a success handler to issue some form of session
		session := sessionStore.New(sessionName)
		session.Values[sessionUserKey] = *githubUser.ID
		session.Save(w)
		http.Redirect(w, req, "/profile", http.StatusFound)
	}
	return ctxh.ContextHandlerFunc(fn)
}

func defaultLoginHandler() http.Handler {

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		log.Println("retrieving github user from context")
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Errorf("Error retrieving github user")))

		}
	}
}
