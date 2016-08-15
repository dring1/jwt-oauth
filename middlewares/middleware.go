package middleware

import (
	"fmt"
	"net/http"
)

type HttpHandler func(http.ResponseWriter, *http.Request) error

func Handlers(handlers ...HttpHandler) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		for _, handler := range handlers {
			fmt.Println(handler)
		}
		return nil
	}
}

func HandlerFuncs(handlers ...http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			handler(w, r)
		}
	})
}
