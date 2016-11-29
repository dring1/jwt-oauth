package routes

import (
	"testing"

	"io/ioutil"
	"net/http"

	"github.com/stretchr/testify/assert"
)

func TestMockServer(t *testing.T) {
	client, mux, server := MockServer()
	defer server.Close()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("hello!"))
	})
	mux.Handle("/test", handler)
	resp, err := client.Get(server.URL + "/test")
	assert.Nil(t, err)
	bs, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "hello!", string(bs))
}

func TestLoginRoute(t *testing.T) {
	client, mux, server := MockServer()
	defer server.Close()

	loginRoute := &GithubLoginRoute{
		Route: Route{
			Path:    "/github/login",
			Methods: []string{Get},
		},
		ClientID:     "TESTID",
		ClientSecret: "TESTSECRET",
	}
	r, err := loginRoute.CompileRoute()

	mux.Handle(r.Path, r.Handler)
	// &GithubCallbackRoute{Route: Route{Path: "/github/callback", Methods: []string{Get}}, ClientID: gitHubClientID, ClientSecret: gitHubClientSecret}
}
