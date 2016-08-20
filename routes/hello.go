package routes

import (
	"net/http"

	"github.com/dring1/jwt-oauth/controllers"
	"github.com/dring1/jwt-oauth/middleware"
)

type HelloRoute struct{}

func (h *HelloRoute) GenHttpHandlers() ([]*R, error) {
	return []*R{
			&R{
				Path:    "/test/hello",
				Handler: middleware.Handlers(http.HandlerFunc(controllers.HelloController)),
				Methods: []string{"GET"},
			},
		},
		nil
}
