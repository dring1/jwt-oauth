package controllers

import (
	"net/http"

	"github.com/dring1/jwt-oauth/cache"
	"github.com/dring1/jwt-oauth/database"
)

// func HelloController(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(201)
// 	w.Write([]byte("<h1>Hello, world!</h1>"))
// }
//
type HelloController struct {
	C
	DatabaseService *database.DatabaseService
	CacheService    *cache.CacheService
}

func (h *HelloController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(201)
	w.Write([]byte("<h1>Hello, world!</h1>"))
}

func (h *HelloController) Route() string {
	return h.C.Route
}

func (h *HelloController) Methods() []string {
	return h.C.Methods
}
