package integration_tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/dring1/jwt-oauth/config"
	"github.com/dring1/jwt-oauth/services"
	"github.com/dring1/jwt-oauth/token"
	"github.com/stretchr/testify/assert"
)

var app *TestApp

func TestNewApp(t *testing.T) {
	c, err := config.New()
	assert.Nil(t, err)
	svcs, err := services.New(c)
	app = NewTestApp(c, svcs)
}

func TestLoginRoute(t *testing.T) {
	authResp := AuthResp{}
	resp, err := app.Client.Get(app.Server.URL + "/mock/github/login")
	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	assert.Nil(t, err)
	assert.NotEmpty(t, authResp.Token)
	assert.Equal(t, "user@acme.com", authResp.Email)
	app.Token = authResp.Token
}

func TestProtectedRouteWithToken(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/test", app.Server.URL), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", app.Token))
	resp, err := app.Client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestProtectedRouteWithoutToken(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/test", app.Server.URL), nil)
	resp, err := app.Client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestProtectedRouteWithInvalidToken(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/test", app.Server.URL), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "Some garbage string"))
	resp, err := app.Client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestProtectedRouteWithExpiredToken(t *testing.T) {
	TtlSeconds := 1
	tokenService, _ := token.NewService(app.Config.PrivateKey, app.Config.PublicKey, TtlSeconds, app.Config.JWTExpirationDelta, app.Config.JwtIss, app.Config.JwtSub, app.Services.CacheService)
	token, err := tokenService.NewToken("vandelay@industries.com")
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	time.Sleep(2 * time.Second)
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/test", app.Server.URL), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := app.Client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestRefreshTokenRoute(t *testing.T) {
	t.Skip()
}
