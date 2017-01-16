package middleware

import "github.com/dring1/jwt-oauth/services"

type MiddlewareKey int
type MiddlewareMap map[MiddlewareKey]Middleware

const (
	ValidateMiddleware = iota
	LoggingMiddleware
	JsonResponseMiddleware
)

// TODO also take config
func New(svcs *services.Services) (MiddlewareMap, error) {
	m := make(MiddlewareMap)
	m[ValidateMiddleware] = NewTokenValidationMiddleware(svcs.TokenService)
	//m[LoggingMiddleware] = NewApacheLoggingHandler()
	return m, nil
}
