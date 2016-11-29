package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/net/context"

	"github.com/dghubble/ctxh"
	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"

	"fmt"

	oauth2Login "github.com/dghubble/gologin/oauth2"
	"github.com/dghubble/sessions"
	s "github.com/dring1/jwt-oauth/app/sessions"
	"github.com/dring1/jwt-oauth/app/users"
	"github.com/dring1/jwt-oauth/lib/errors"
	githubClient "github.com/google/go-github/github"
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
	ClientID       string
	ClientSecret   string
	RedirectURL    string
	UserService    users.Service
	SessionService s.Service
}

func (ghr *GithubLoginRoute) CompileRoute() (*Route, error) {
	config := &oauth2.Config{
		ClientID:     ghr.ClientID,
		ClientSecret: ghr.ClientSecret,
		RedirectURL:  ghr.RedirectURL,
		Endpoint:     githubOAuth2.Endpoint,
		Scopes:       []string{"user:email"},
	}

	stateConfig := gologin.DebugOnlyCookieConfig

	ghr.Handler = ctxh.NewHandler(github.StateHandler(stateConfig, github.LoginHandler(config, nil)))
	return &ghr.Route, nil
}

type GithubCallbackRoute struct {
	Route
	ClientID       string
	ClientSecret   string
	RedirectURL    string
	UserService    users.Service `service:"userService"`
	SessionService s.Service     `service:"sessionService"`
}

func (gcr *GithubCallbackRoute) CompileRoute() (*Route, error) {
	config := &oauth2.Config{
		ClientID:     gcr.ClientID,
		ClientSecret: gcr.ClientSecret,
		RedirectURL:  gcr.RedirectURL,
		Endpoint:     githubOAuth2.Endpoint,
		Scopes:       []string{"user:email"},
	}

	if gcr.Handler == nil {
		stateConfig := gologin.DebugOnlyCookieConfig
		gcr.Handler = ctxh.NewHandler(github.StateHandler(stateConfig, github.CallbackHandler(config, gcr.defaultLoginHandler(), nil)))

	}
	return &gcr.Route, nil
}

func (gcr *GithubCallbackRoute) defaultLoginHandler() ctxh.ContextHandler {

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		log.Println("retrieving github user from context")
		githubUser, err := github.UserFromContext(ctx)
		t, err := oauth2Login.TokenFromContext(ctx)
		validateGithubUser(t.AccessToken, githubUser.Email)
		// fmt.Println(t.)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error retrieving github user"))

			return
		}
		// log.Println(githubUser)
		err = gcr.UserService.Authenticate(*githubUser.Email)
		if err != nil {
			// Failed to Authenticate
			w.WriteHeader(http.StatusUnauthorized)
			errors.ErrorHandler(w, r)
			return
		}

		token, err := gcr.SessionService.NewSession(*githubUser.Email)
		if err != nil {
			w.WriteHeader(500)
			errors.ErrorHandler(w, r)
			return
		}
		err = json.NewEncoder(w).Encode(token)
		if err != nil {
			w.WriteHeader(500)
			errors.ErrorHandler(w, r)
			return
		}
	}

	return ctxh.ContextHandlerFunc(handler)
}

func validateGithubUser(token string, expectedEmail *string) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := githubClient.NewClient(tc)

	// list all repositories for the authenticated user
	emails, _, err := client.Users.ListEmails(nil)
	if err != nil {
		fmt.Println(err)
	}

	for _, email := range emails {
		if *email.Email == *expectedEmail {
			fmt.Println("GREAT SUCCESS")
		}
	}
}
