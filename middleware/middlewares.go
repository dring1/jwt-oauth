package middleware

import "github.com/dring1/jwt-oauth/services"

type MiddlewareKey int
type MiddlewareMap map[MiddlewareKey]Middleware

const (
	ValidateMiddleware = iota
	LoggingMiddleware
)

// TODO also take config
func New(svcs *services.Services) (MiddlewareMap, error) {
	m := make(MiddlewareMap)
	m[ValidateMiddleware] = NewTokenValidationMiddleware(svcs.TokenService)
<<<<<<< HEAD
	m[LoggingMiddleware] = NewApacheLoggingHandler()
=======
	//m[LoggingMiddleware] = NewApacheLoggingHandler(svcs.)
>>>>>>> 4fab027b6560c49e20448e9d6d0f6f33d55d5287
	return m, nil
}
