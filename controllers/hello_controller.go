package controllers

import (
	"net/http"

	"github.com/dring1/jwt-oauth/cache"
	"github.com/dring1/jwt-oauth/database"
)

func HelloController(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(201)
	w.Write([]byte("<h1>Hello, world!</h1>"))
}

type HellController struct {
	C
	DatabaseService *database.DatabaseService
	CacheService    *cache.CacheService
}

func (h *HellController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(201)
	w.Write([]byte("<h1>Hello, world!</h1>"))
}

func (h *HellController) Route() string {
	return "/Word"
}

func (h *HellController) Methods() []string {
	return []string{"GET"}
}
