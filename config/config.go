package config

import (
	"io"
	"time"
)

type Cfg struct {
	JWTExpirationDelta int
	PrivateKey         []byte
	PublicKey          []byte
	Port               int
	GitHubClientID     string
	GitHubClientSecret string
	OauthRedirectURL   string
	LoggingEndpoint    io.Writer
	LogLevel           string
	RedisEndpoint      string
	JwtTTL             int
	JwtIss             string
	JwtSub             string
}

func NewConfig(opts ...func(*Cfg) error) (*Cfg, error) {
	c := &Cfg{
		JWTExpirationDelta: (int)(time.Hour.Hours()),
		PrivateKey:         make([]byte, 10),
		PublicKey:          make([]byte, 10),
		Port:               8080,
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}
