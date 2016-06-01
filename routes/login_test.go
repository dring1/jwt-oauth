package routes

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dghubble/ctxh"
	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"
	oauth2Login "github.com/dghubble/gologin/oauth2"
	"github.com/dghubble/gologin/testutils"
	"github.com/dring1/jwt-oauth/controllers"
	"github.com/dring1/jwt-oauth/models"
	gh "github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func TestMockHandler(t *testing.T) {
	m := mux.NewRouter()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("hello!"))
	})
	m.Handle("/", handler)
	c, _, s := MockServer(m)
	defer s.Close()
	resp, _ := c.Get(s.URL + "/")
	bs, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "hello!", string(bs))
}

func TestGithubMockHandler(t *testing.T) {
	config := &oauth2.Config{}
	jsonData := `{"id": 917408, "name": "Alyssa Hacker"}`
	expectedUser := &gh.User{ID: gh.Int(917408), Name: gh.String("Alyssa Hacker")}
	proxyClient, server := newGithubTestServer(jsonData)
	defer server.Close()
	// oauth2 Client will use the proxy client's base Transport
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, proxyClient)
	anyToken := &oauth2.Token{AccessToken: "any-token"}
	ctx = oauth2Login.WithToken(ctx, anyToken)

	success := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		githubUser, err := github.UserFromContext(ctx)
		assert.Nil(t, err)
		assert.Equal(t, expectedUser, githubUser)
		fmt.Fprintf(w, "success handler called")
	}
	failure := testutils.AssertFailureNotCalled(t)

	// GithubHandler assert that:
	// - Token is read from the ctx and passed to the Github API
	// - github User is obtained from the Github API
	// - success handler is called
	// - github User is added to the ctx of the success handler
	githubHandler := githubHandler(config, ctxh.ContextHandlerFunc(success), failure)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)
	githubHandler.ServeHTTP(ctx, w, req)
	assert.Equal(t, "success handler called", w.Body.String())
}

func TestSuccessfulLogin(t *testing.T) {
	config := &oauth2.Config{}
	jsonData := `{"id": 917408, "name": "Alyssa Hacker"}`
	expectedUser := &gh.User{ID: gh.Int(917408), Name: gh.String("Alyssa Hacker")}
	proxyClient, server := newGithubTestServer(jsonData)
	defer server.Close()
	// oauth2 Client will use the proxy client's base Transport
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, proxyClient)
	anyToken := &oauth2.Token{AccessToken: "any-token"}
	ctx = oauth2Login.WithToken(ctx, anyToken)
	success := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		githubUser, err := github.UserFromContext(ctx)
		assert.Nil(t, err)
		assert.Equal(t, expectedUser, githubUser)
		fmt.Fprintf(w, "success handler called")
	}
	failure := testutils.AssertFailureNotCalled(t)

	// GithubHandler assert that:
	// - Token is read from the ctx and passed to the Github API
	// - github User is obtained from the Github API
	// - success handler is called
	// - github User is added to the ctx of the success handler
	githubHandler := githubHandler(config, ctxh.ContextHandlerFunc(success), failure)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		githubHandler.ServeHTTP(ctx, w, r)
	})
	m := mux.NewRouter()
	m = LoginRoute(m, handler, nil, config)
	c, _, s := MockServer(m)
	defer s.Close()
	resp, err := c.Get(s.URL + "/github/login")
	assert.Nil(t, err)
	bs, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "success handler called", string(bs))
}

func TestSuccessfulLoginWithJWT(t *testing.T) {
	config := &oauth2.Config{}
	jsonData := `{"id": 917408, "name": "Alyssa Hacker", "email": "user.haxor@yahoo.com"}`
	expectedUser := &gh.User{ID: gh.Int(917408), Name: gh.String("Alyssa Hacker"), Email: gh.String("uber.haxor@yahoo.com")}
	proxyClient, server := newGithubTestServer(jsonData)
	defer server.Close()
	// oauth2 Client will use the proxy client's base Transport
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, proxyClient)
	anyToken := &oauth2.Token{AccessToken: "any-token"}
	ctx = oauth2Login.WithToken(ctx, anyToken)
	// success := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	// 	githubUser, err := github.UserFromContext(ctx)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, expectedUser, githubUser)
	// 	fmt.Fprintf(w, "success handler called")
	// }
	success := controllers.Login(func(githubUser *models.User) {
		assert.Equal(t, expectedUser, githubUser)
	})
	failure := testutils.AssertFailureNotCalled(t)

	// GithubHandler assert that:
	// - Token is read from the ctx and passed to the Github API
	// - github User is obtained from the Github API
	// - success handler is called
	// - github User is added to the ctx of the success handler
	githubHandler := githubHandler(config, success, failure)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		githubHandler.ServeHTTP(ctx, w, r)
	})
	m := mux.NewRouter()
	m = LoginRoute(m, handler, nil, config)
	c, _, s := MockServer(m)
	defer s.Close()
	resp, err := c.Get(s.URL + "/github/login")
	assert.Nil(t, err)
	bs, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "success handler called", string(bs))
}

func githubHandler(config *oauth2.Config, success, failure ctxh.ContextHandler) ctxh.ContextHandler {
	if failure == nil {
		failure = gologin.DefaultFailureHandler
	}
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		token, err := oauth2Login.TokenFromContext(ctx)
		if err != nil {
			log.Println("Err", err)
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(ctx, w, req)
			return
		}
		httpClient := config.Client(ctx, token)
		githubClient := gh.NewClient(httpClient)
		user, resp, err := githubClient.Users.Get("")
		err = validateResponse(user, resp, err)
		if err != nil {
			log.Println(err)
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(ctx, w, req)
			return
		}
		ctx = github.WithUser(ctx, user)
		success.ServeHTTP(ctx, w, req)
	}
	return ctxh.ContextHandlerFunc(fn)
}

// validateResponse returns an error if the given Github user, raw
// http.Response, or error are unexpected. Returns nil if they are valid.
func validateResponse(user *gh.User, resp *gh.Response, err error) error {
	if err != nil || resp.StatusCode != http.StatusOK {
		return github.ErrUnableToGetGithubUser
	}
	if user == nil || user.ID == nil {
		return github.ErrUnableToGetGithubUser
	}
	return nil
}

type key int

const (
	userKey key = iota
)

// responds with the given json data. The caller must close the server.
// func newGithubTestServer(jsonData string, config *oauth2.Config) (*http.Client, *httptest.Server) {
// 	router := mux.NewRouter()
// 	f := func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		fmt.Fprintf(w, jsonData)
// 	}
// 	router = LoginRoute(router, http.HandlerFunc(f), nil, config)
// 	client, _, server := MockServer(router)
// 	return client, server
// }

func MockServer(router *mux.Router) (*http.Client, *mux.Router, *httptest.Server) {
	server := httptest.NewServer(router)
	transport := &testutils.RewriteTransport{&http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}}
	client := &http.Client{Transport: transport}
	return client, router, server
}

func newGithubTestServer(jsonData string) (*http.Client, *httptest.Server) {
	client, mux, server := testutils.TestServer()
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, jsonData)
	})
	return client, server
}
