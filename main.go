package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dring1/jwt-oauth/app/sessions"
	"github.com/dring1/jwt-oauth/app/users"
	"github.com/dring1/jwt-oauth/cache"
	"github.com/dring1/jwt-oauth/middleware"
	"github.com/dring1/jwt-oauth/routes"
	"github.com/dring1/jwt-oauth/token"
)

func main() {
	c, err := NewAppConfig()
	if err != nil {
		log.Fatalln(err)
	}
	// Init services
	// db, _ := database.NewService()
	ch, _ := cache.NewService(c.RedisEndpoint)
	tokenService, _ := token.NewService(c.PrivateKey, c.PublicKey, c.JwtTTL, c.JWTExpirationDelta, c.JwtIss, c.JwtSub)
	us, _ := users.NewService()
	ss, _ := sessions.NewService(tokenService, ch)
	// Init controllers
	// ctrls := controllers.New(db, ch, jwtService, us, ss)
	routeServices := map[string]interface{}{
		"userService":    us,
		"sessionService": ss,
	}

	tokenValidationMiddlware := middleware.NewTokenValidationMiddleware(tokenService)
	middlewares := map[string]middleware.Middleware{
		"VALIDATION": tokenValidationMiddlware,
	}
	// Init router
	router := routes.New(c.GitHubClientID, c.GitHubClientSecret, c.OauthRedirectURL, routeServices, middlewares)

	// Apply middlewares
	globalMiddlewares := []middleware.Middleware{
		middleware.NewApacheLoggingHandler(c.LoggingEndpoint),
	}
	globalMiddlewares = append(globalMiddlewares, middleware.DefaultMiddleWare()...)

	log.Printf("Serving on port :%d", c.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", c.Port), middleware.Handlers(router, globalMiddlewares...))
	log.Fatal(err)
}
