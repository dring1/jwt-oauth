package routes

import (
	"context"
	"log"
	"net/http"

	"fmt"

	"github.com/dring1/jwt-oauth/app/users"
	"github.com/dring1/jwt-oauth/lib/contextkeys"
	"github.com/dring1/jwt-oauth/lib/errors"
	"github.com/dring1/jwt-oauth/token"
	githubClient "github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githubOAuth2 "golang.org/x/oauth2/github"
)

const (
	sessionName    = "example-github-app"
	sessionSecret  = "example cookie signing secret"
	sessionUserKey = "githubID"
)

var oauthStateString = "ThisIsASecret"

type GithubLoginRoute struct {
	Route
	ClientID     string
	ClientSecret string
	RedirectURL  string
	UserService  users.Service
	LoginHandler http.Handler
	Config       *oauth2.Config
}

func (ghr *GithubLoginRoute) CompileRoute(Responder) (*Route, error) {
	ghr.Config = &oauth2.Config{
		ClientID:     ghr.ClientID,
		ClientSecret: ghr.ClientSecret,
		RedirectURL:  ghr.RedirectURL,
		Endpoint:     githubOAuth2.Endpoint,
		Scopes:       []string{"user:email"},
	}

	url := ghr.Config.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	ghr.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})

	return &ghr.Route, nil
}

type GithubCallbackRoute struct {
	Route
	ClientID     string
	ClientSecret string
	RedirectURL  string
	UserService  users.Service `service:"UserService"`
	TokenService token.Service `service:"TokenService"`
	Config       *oauth2.Config
}

func (gcr *GithubCallbackRoute) CompileRoute(responder Responder) (*Route, error) {
	gcr.Config = &oauth2.Config{
		ClientID:     gcr.ClientID,
		ClientSecret: gcr.ClientSecret,
		RedirectURL:  gcr.RedirectURL,
		Endpoint:     githubOAuth2.Endpoint,
		Scopes:       []string{"user:email"},
	}

	gcr.Handler = gcr.NewHandleGitHubCallback(responder)
	return &gcr.Route, nil
}

func validateGithubUser(token string, expectedEmail *string) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := githubClient.NewClient(tc)

	// list all repositories for the authenticated user
	emails, _, err := client.Users.ListEmails(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
	}

	for _, email := range emails {
		if *email.Email == *expectedEmail {
			fmt.Println("GREAT SUCCESS")
		}
	}
}

func (gcr *GithubCallbackRoute) NewHandleGitHubCallback(responder Responder) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		state := r.FormValue("state")
		if state != oauthStateString {
			//error out here
			log.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		code := r.FormValue("code")
		token, err := gcr.Config.Exchange(oauth2.NoContext, code)
		if err != nil {
			//error out here
			log.Printf("oauthConf.Exchange() failed with '%s'\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		oauthClient := gcr.Config.Client(oauth2.NoContext, token)
		client := githubClient.NewClient(oauthClient)
		user, _, err := client.Users.Get(r.Context(), "")
		if err != nil {
			//error out here
			log.Printf("client.Users.Get() faled with '%s'\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		err = gcr.UserService.Authenticate(*user.Email)
		if err != nil {
			// Failed to Authenticate
			w.WriteHeader(http.StatusUnauthorized)
			errors.ErrorHandler(w, r)
			return
		}

		//jwtToken, err := gcr.SessionService.NewSession(*user.Email)
		jwtToken, err := gcr.TokenService.NewToken(*user.Email)
		if err != nil {
			w.WriteHeader(500)
			errors.ErrorHandler(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), contextkeys.Value, jwtToken)
		r = r.WithContext(ctx)
		log.Println("calling responder")
		responder.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(fn)
}
