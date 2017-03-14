package routes

import "net/http"

type TestRoute struct {
	Route
}

func (r *TestRoute) CompileRoute() (*Route, error) {

	fn := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte("Great success"))
	}
	r.Handler = http.HandlerFunc(fn)
	return &r.Route, nil
}
