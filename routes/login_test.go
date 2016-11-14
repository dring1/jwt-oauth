package routes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
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

}
