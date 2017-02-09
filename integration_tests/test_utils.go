package integration_tests

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/dring1/jwt-oauth/config"
	"github.com/dring1/jwt-oauth/lib/contextkeys"
	"github.com/dring1/jwt-oauth/middleware"
	"github.com/dring1/jwt-oauth/routes"
	"github.com/dring1/jwt-oauth/services"
)

type RewriteTransport struct {
	Transport http.RoundTripper
}
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

func MockServer(mux http.Handler) (*http.Client, http.Handler, *httptest.Server) {
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
	authRoutes := mockAuthRoutes(svcs)
	rs = append(rs, authRoutes...)
	m, err := routes.NewRouter(rs)
	// TODO: move to default middlewares ?
	defaultMiddlewares := middleware.DefaultMiddleWare(config)
	handler := middleware.Handlers(m, defaultMiddlewares...)
	if err != nil {
		log.Fatal(err)
	}
	client, _, server := MockServer(handler)
	return &TestApp{
		Config: config,
		Client: client,
		Server: server,
		//Router:      mux,
		Services:    svcs,
		Middlewares: middlewares,
	}
}

func mockAuthRoutes(svcs *services.Services) []*routes.Route {
	responder := routes.NewResponder()
	loginRoute := &routes.GithubLoginRoute{
		Route: routes.Route{
			Path:    "/mock/github/login",
			Methods: []string{routes.Get},
		},
		ClientID:     "TESTID",
		ClientSecret: "TESTSECRET",
		LoginHandler: nil,
	}
	r, _ := loginRoute.CompileRoute(responder)
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
	cr, _ := callBackRoute.CompileRoute(responder)
	cr.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		jwtToken, err := svcs.TokenService.NewToken("user@acme.com")
		if err != nil {
			w.WriteHeader(500)
			return
		}
		ctx := context.WithValue(r.Context(), contextkeys.Value, jwtToken)
		r = r.WithContext(ctx)
		routes.NewResponder().ServeHTTP(w, r)
	})

	return []*routes.Route{r, cr}
}
