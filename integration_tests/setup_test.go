package integration_tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/dring1/jwt-oauth/app/sessions"
	"github.com/dring1/jwt-oauth/app/users"
	"github.com/dring1/jwt-oauth/cache"
	"github.com/dring1/jwt-oauth/config"
	"github.com/dring1/jwt-oauth/middleware"
	"github.com/dring1/jwt-oauth/routes"
	"github.com/dring1/jwt-oauth/token"
	"github.com/stretchr/testify/assert"
)

var app *TestApp
var services *TestServices

func TestNewApp(t *testing.T) {
	c, err := config.New()
	assert.Nil(t, err)
	ch, err := cache.NewService(c.RedisEndpoint)
	assert.Nil(t, err)
	tokenService, _ := token.NewService(c.PrivateKey, c.PublicKey, c.JwtTTL, c.JWTExpirationDelta, c.JwtIss, c.JwtSub, ch)
	assert.Nil(t, err)
	us, err := users.NewService()
	assert.Nil(t, err)
	ss, err := sessions.NewService(tokenService, ch)
	assert.Nil(t, err)
	services = &TestServices{
		cacheService:   ch,
		tokenService:   tokenService,
		userService:    us,
		sessionService: ss,
	}
	app = NewTestApp(services)
}

func TestLoginRoute(t *testing.T) {
	authResp := AuthResp{}
	resp, err := app.Client.Get(app.Server.URL + "/github/login")
	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	assert.Nil(t, err)
	assert.NotEmpty(t, authResp.Token)
	assert.Equal(t, "user@acme.com", authResp.Email)
	app.Token = authResp.Token
}

func TestProtectedRouteWithToken(t *testing.T) {
	route := &routes.TestRoute{Route: routes.Route{Path: "/test", Methods: []string{"Get"}, Middlewares: []middleware.Middleware{app.Middlewares["VALIDATION"]}}}
	r, _ := route.CompileRoute()
	assert.Equal(t, 1, len(r.Middlewares))
	handler := middleware.Handlers(r.Handler, r.Middlewares...)
	app.Mux.Handle(r.Path, handler)
	req, _ := http.NewRequest("GET", app.Server.URL+"/test", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", app.Token))
	resp, err := app.Client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)

}

func TestProtectedRouteWithoutToken(t *testing.T) {
	//route := &routes.TestRoute{Route: routes.Route{Path: "/test", Methods: []string{"Get"}, Middlewares: []middleware.Middleware{app.Middlewares["VALIDATION"]}}}
	//route.CompileRoute()
	//app.Mux.Handle(route.Path, route.Handler)
	req, _ := http.NewRequest("GET", app.Server.URL+"/test", nil)
	resp, err := app.Client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, 401, resp.StatusCode)

}
