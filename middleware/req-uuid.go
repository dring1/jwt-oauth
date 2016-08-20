package middleware

import (
	"log"
	"net/http"
)

func AddUUID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Applying middleware")
		// next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
