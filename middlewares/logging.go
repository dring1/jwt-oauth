package middlewares

import (
	"log"
	"net/http"
	"time"
)

func LoggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %q %v\n", r.Method, r.URL.String(), r.UserAgent(), t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}
