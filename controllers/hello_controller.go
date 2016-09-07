package controllers

import (
	"github.com/dring1/jwt-oauth/cache"
	"github.com/dring1/jwt-oauth/database"
)

// func HelloController(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(201)
// 	w.Write([]byte("<h1>Hello, world!</h1>"))
// }
//
type HelloController struct {
	// C
	DatabaseService *database.DatabaseService
	CacheService    *cache.CacheService
}

// func (h *HelloController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(201)
// 	w.Write([]byte("<h1>Hello, world!</h1>"))
// }
//
// func (h *HelloController) GetRoute() string {
// 	return h.C.Route
// }
//
// func (h *HelloController) GetMethods() []string {
// 	return h.C.Methods
// }

func (c *HelloController) Create(obj interface{}) error {
	return nil
}
func (c *HelloController) Read(obj interface{}) error {
	return nil
}
func (c *HelloController) Update(obj interface{}) error {
	return nil
}
func (c *HelloController) Delete(obj interface{}) error {
	return nil
}
