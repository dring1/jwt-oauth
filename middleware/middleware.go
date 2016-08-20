package middleware

import "net/http"

func HandlerFuncs(handlers ...http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			handler(w, r)
		}
	})
}

func Handlers(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, mdlware := range middlewares {
		handler = mdlware(handler)
	}
	return handler
}

type Middleware func(http.Handler) http.Handler

func DefaultMiddleWare() []Middleware {
	return []Middleware{RecoverHandler}
}
