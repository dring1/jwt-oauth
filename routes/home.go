package routes

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func HomeRoute(r *mux.Router) *mux.Router {
	fp := "static"
	_, err := filepath.Abs(fp)
	if err != nil {
		log.Fatal(err)
	}
	r.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(fp))))
	return r
}
