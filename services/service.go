package services

import (
	"github.com/dring1/jwt-oauth/app/users"
	"github.com/dring1/jwt-oauth/cache"
	"github.com/dring1/jwt-oauth/config"
	"github.com/dring1/jwt-oauth/database"
	jsonresponder "github.com/dring1/jwt-oauth/jsonResponder"
	"github.com/dring1/jwt-oauth/logger"
	"github.com/dring1/jwt-oauth/token"
)

type Services struct {
	UserService   users.Service
	TokenService  token.Service
	Cache         *cache.Service
	Logger        logger.Service
	JsonResponder jsonresponder.Service
	Database      *database.Service
}

func New(c *config.Cfg) (*Services, error) {
	dbConfig := &database.Config{
		Host:     c.DbHost,
		Port:     c.DbPort,
		User:     c.DbUser,
		Password: c.DbPassword,
		DbName:   c.DbName,
		SSL:      c.DbSSL,
	}
	dbService, err := database.NewService(dbConfig)
	if err != nil {
		return nil, err
	}
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

	jsonResponder := jsonresponder.NewJsonResponder()

	services := &Services{
		UserService:   userService,
		TokenService:  tokenService,
		Cache:         cacheService,
		Logger:        loggerService,
		JsonResponder: jsonResponder,
		Database:      dbService,
	}
	return services, nil
}
