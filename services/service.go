package services

import (
	"github.com/dring1/jwt-oauth/app/users"
	"github.com/dring1/jwt-oauth/cache"
	"github.com/dring1/jwt-oauth/config"
	"github.com/dring1/jwt-oauth/logger"
	"github.com/dring1/jwt-oauth/token"
)

type Services struct {
	UserService   users.Service
	TokenService  token.Service
	CacheService  *cache.Service
	LoggerService logger.Service
}

func New(c *config.Cfg) (*Services, error) {
	cacheService, err := cache.NewService(c.RedisEndpoint)
	if err != nil {
		return nil, err
	}
	tokenService, err := token.NewService(c.PrivateKey, c.PublicKey, c.JwtTTL, c.JWTExpirationDelta, c.JwtIss, c.JwtSub, cacheService)
	if err != nil {
		return nil, err
	}
	userService, err := users.NewService()
	if err != nil {
		return nil, err
	}
	loggerConfig := &logger.Config{
		Level:    c.LogLevel,
		Endpoint: c.LoggingEndpoint,
	}
	loggerService := logger.NewLoggerService(loggerConfig)

	services := &Services{
		UserService:   userService,
		TokenService:  tokenService,
		CacheService:  cacheService,
		LoggerService: loggerService,
	}
	return services, nil
}
