package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	middleware "github.com/dring1/jwt-oauth/middlewares"
	"github.com/stretchr/testify/assert"
)

func mockMiddleWare() []http.HandlerFunc {
	mdw := []http.HandlerFunc{}
	for index := 0; index < 10; index++ {
		h := (func(i int) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add(fmt.Sprintf("MIDDLEWARE-%d", i), strconv.Itoa(i))
			}

		})(index)
		mdw = append(mdw, h)
	}
	return mdw
}
func TestApplyManyMiddleWares(t *testing.T) {
	middlewareSlice := mockMiddleWare()
	assert.Equal(t, len(middlewareSlice), 10)

	handler := middleware.HandlerFuncs(middlewareSlice...)

	req, err := http.NewRequest("GET", "localhost:8080", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	handler(w, req)
	x := w.Header().Get("MIDDLEWARE-5")
	assert.Equal(t, "5", x)
}
