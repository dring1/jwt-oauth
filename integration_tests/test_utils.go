package integration_tests

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/dring1/jwt-oauth/config"
	"github.com/dring1/jwt-oauth/middleware"
	"github.com/dring1/jwt-oauth/routes"
	"github.com/dring1/jwt-oauth/services"
)

type RewriteTransport struct {
	Transport http.RoundTripper
}

//type TestServices struct {
//    cacheService   *cache.Service
//    userService    users.Service
//    tokenService   token.Service
//    sessionService sessions.Service
//}

type AuthResp struct {
	Token string `json:"Value"`
	Email string `json:"Email"`
}

type TestApp struct {
	Config      *config.Cfg
	Client      *http.Client
	Router      *routes.Router
	Server      *httptest.Server
	Token       string
	Services    *services.Services
	Middlewares middleware.MiddlewareMap
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

func MockServer(mux *routes.Router) (*http.Client, *routes.Router, *httptest.Server) {
	server := httptest.NewServer(mux)
	transport := &RewriteTransport{&http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}}
	client := &http.Client{Transport: transport}
	return client, mux, server
}

func NewTestApp(config *config.Cfg, svcs *services.Services) *TestApp {
	middlewares, _ := middleware.New(svcs)
	rs, _ := routes.NewRoutes(&routes.Config{
		Middlewares: middlewares,
		Services:    svcs,
	})
	m, err := routes.NewRouter(rs)
	if err != nil {
		log.Fatal(err)
	}
	client, mux, server := MockServer(m)

	loginRoute := &routes.GithubLoginRoute{
		Route: routes.Route{
			Path:    "/mock/github/login",
			Methods: []string{routes.Get},
		},
		ClientID:     "TESTID",
		ClientSecret: "TESTSECRET",
		LoginHandler: nil,
	}
	r, _ := loginRoute.CompileRoute()
	r.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/mock/github/callback", 301)
	})
	callBackRoute := &routes.GithubLoginRoute{
		Route: routes.Route{
			Path:    "/mock/github/callback",
			Methods: []string{routes.Get},
		},
		ClientID:     "TESTID",
		ClientSecret: "TESTSECRET",
		LoginHandler: nil,
	}
	cr, _ := callBackRoute.CompileRoute()
	cr.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		jwtToken, err := svcs.SessionService.NewSession("user@acme.com")
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
		Config:      config,
		Client:      client,
		Server:      server,
		Router:      mux,
		Services:    svcs,
		Middlewares: middlewares,
	}
	//router := routes.New(c.GitHubClientID, c.GitHubClientSecret, c.OauthRedirectURL, routeServices, middlewares)

	// Apply middlewares
	//globalMiddlewares := []middleware.Middleware{
	//    middleware.NewApacheLoggingHandler(c.LoggingEndpoint),
	//}
	//globalMiddlewares = append(globalMiddlewares, middleware.DefaultMiddleWare()...)

	//log.Printf("Serving on port :%d", c.Port)
	//err = http.ListenAndServe(fmt.Sprintf(":%d", c.Port), middleware.Handlers(router, globalMiddlewares...))
}

//func MockMiddlewares(services map[string]interface{}) map[string]middleware.Middleware {
//    tokenValidationMiddleware := middleware.NewTokenValidationMiddleware(services["tokenService"])
//    middlewares := map[string]middleware.Middleware{
//        "VALIDATION": tokenValidationMiddleware,
//    }
//    return middlewares
//}
