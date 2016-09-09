package routes

import (
	"fmt"
	"net/http"
)

type UserProfileRoute struct {
	Route
}

func (upr *UserProfileRoute) NewHandler() (*R, error) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, `<p>You are logged in!</p><form action="/logout" method="post"><input type="submit" value="Logout"></form>`)
	}
	return &R{
		Path:    upr.Path,
		Methods: upr.Methods,
		Handler: http.HandlerFunc(handler),
	}, nil
}
