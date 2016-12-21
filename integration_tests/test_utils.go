package integration_tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/dring1/jwt-oauth/app/sessions"
	"github.com/dring1/jwt-oauth/app/users"
	"github.com/dring1/jwt-oauth/cache"
	"github.com/dring1/jwt-oauth/middleware"
	"github.com/dring1/jwt-oauth/routes"
	"github.com/dring1/jwt-oauth/token"
)

type RewriteTransport struct {
	Transport http.RoundTripper
}

type TestServices struct {
	cacheService   *cache.Service
	userService    users.Service
	tokenService   token.Service
	sessionService sessions.Service
}

type AuthResp struct {
	Token string `json:"Value"`
	Email string `json:"Email"`
}

type TestApp struct {
	Client      *http.Client
	Mux         *http.ServeMux
	Server      *httptest.Server
	Token       string
	Services    *TestServices
	Middlewares map[string]middleware.Middleware
}

// RoundTrip rewrites the request scheme to http and calls through to the
// composed RoundTripper or if it is nil, to the http.DefaultTransport.
func (t *RewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	if t.Transport == nil {
		return http.DefaultTransport.RoundTrip(req)
	}
	return t.Transport.RoundTrip(req)
}

func MockServer() (*http.Client, *http.ServeMux, *httptest.Server) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	transport := &RewriteTransport{&http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}}
	client := &http.Client{Transport: transport}
	return client, mux, server
}

func NewTestApp(services *TestServices) *TestApp {
	client, mux, server := MockServer()
	loginRoute := &routes.GithubLoginRoute{
		Route: routes.Route{
			Path:    "/github/login",
			Methods: []string{routes.Get},
		},
		ClientID:     "TESTID",
		ClientSecret: "TESTSECRET",
		LoginHandler: nil,
	}
	r, _ := loginRoute.CompileRoute()
	r.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/github/callback", 301)
	})
	callBackRoute := &routes.GithubLoginRoute{
		Route: routes.Route{
			Path:    "/github/callback",
			Methods: []string{routes.Get},
		},
		ClientID:     "TESTID",
		ClientSecret: "TESTSECRET",
		LoginHandler: nil,
	}
	cr, _ := callBackRoute.CompileRoute()
	cr.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		jwtToken, err := services.sessionService.NewSession("user@acme.com")
		if err != nil {
			w.WriteHeader(500)
			return
		}
		err = json.NewEncoder(w).Encode(jwtToken)
		if err != nil {
			w.WriteHeader(500)
			return
		}

	})
	mux.Handle(callBackRoute.Path, callBackRoute.Handler)
	mux.Handle(loginRoute.Path, loginRoute.Handler)
	return &TestApp{
		Client:      client,
		Server:      server,
		Mux:         mux,
		Services:    services,
		Middlewares: MockMiddlewares(services),
	}
}

func MockMiddlewares(services *TestServices) map[string]middleware.Middleware {
	tokenValidationMiddleware := middleware.NewTokenValidationMiddleware(services.tokenService)
	middlewares := map[string]middleware.Middleware{
		"VALIDATION": tokenValidationMiddleware,
	}
	return middlewares
}
