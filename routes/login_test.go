package routes

import (
	"testing"

	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func LoginServer() (*http.Client, *http.ServeMux, *httptest.Server) {
	success := "success"
	client, mux, server := MockServer()
	loginRoute := &GithubLoginRoute{
		Route: Route{
			Path:    "/github/login",
			Methods: []string{Get},
		},
		ClientID:     "TESTID",
		ClientSecret: "TESTSECRET",
		LoginHandler: nil,
	}
	r, _ := loginRoute.CompileRoute()
	r.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/github/callback", 301)
	})
	callBackRoute := &GithubLoginRoute{
		Route: Route{
			Path:    "/github/callback",
			Methods: []string{Get},
		},
		ClientID:     "TESTID",
		ClientSecret: "TESTSECRET",
		LoginHandler: nil,
	}
	cr, _ := callBackRoute.CompileRoute()
	cr.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte(success))
	})
	mux.Handle(callBackRoute.Path, callBackRoute.Handler)
	mux.Handle(loginRoute.Path, loginRoute.Handler)
	return client, mux, server
}
func TestLoginRoute(t *testing.T) {
	client, _, server := LoginServer()
	defer server.Close()

	resp, err := client.Get(server.URL + "/github/callback")
	assert.Nil(t, err)
	b, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "success", string(b))
	assert.Equal(t, 201, resp.StatusCode)
}

func TestLogin(t *testing.T) {

}
