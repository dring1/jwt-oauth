package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
)

type DefaultValFunc func() (interface{}, error)

func New() (*Cfg, error) {
	var PrivateKey *pem.Block
	privateKey := func(c *Cfg) error {
		privateKeyPemBlock, err := getEnvVal("PRIVATE_KEY", func() (interface{}, error) {
			pk, _ := rsa.GenerateKey(rand.Reader, 1024)
			bits := x509.MarshalPKCS1PrivateKey(pk)
			pemBlock := pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: bits,
			}
			PrivateKey = &pemBlock
			return &pemBlock, nil
		})
		if err != nil {
			return err
		}
		c.PrivateKey = privateKeyPemBlock.(*pem.Block).Bytes
		return nil
	}
	publicKey := func(c *Cfg) error {
		getEnvVal("PUBLIC_KEY", func() (interface{}, error) {
			pKey := PrivateKey.Bytes
			privKey, err := x509.ParsePKCS1PrivateKey(pKey)
			if err != nil {
				return nil, err
			}
			pubKey := privKey.PublicKey
			pub, err := x509.MarshalPKIXPublicKey(&pubKey)
			if err != nil {
				return nil, err
			}
			pemBlock := pem.Block{
				Type:  "PUBLIC KEY",
				Bytes: pub,
			}
			return &pemBlock, nil
		})
		return nil
	}

	port := func(c *Cfg) error {
		p, err := getEnvVal("PORT", func() (interface{}, error) {
			return 8080, nil
		})
		if err != nil {
			return err
		}
		c.Port = p.(int)
		return nil
	}

	gitHubClientID := func(c *Cfg) error {
		ghCID, err := getEnvVal("GITHUB_CLIENT_ID", func() (interface{}, error) {
			return nil, errors.Errorf("Did not provide GITHUB_CLIENT_ID")
		})
		if err != nil {
			return err
		}
		c.GitHubClientID = ghCID.(string)
		return nil
	}

	gitHubClientSecret := func(c *Cfg) error {
		ghCS, err := getEnvVal("GITHUB_CLIENT_SECRET", func() (interface{}, error) {
			return nil, errors.Errorf("Did not provide GITHUB_CLIENT_SECRET")
		})
		if err != nil {
			return err
		}
		c.GitHubClientSecret = ghCS.(string)
		return nil
	}

	oauthRedirectURL := func(c *Cfg) error {
		rdURL, err := getEnvVal("OAUTH_REDIRECT_URL", func() (interface{}, error) {
			return fmt.Sprintf("http://localhost:%d/github/callback", c.Port), nil
		})
		if err != nil {
			return err
		}
		c.OauthRedirectURL = rdURL.(string)
		return nil
	}
	loggingEndpoint := func(c *Cfg) error {
		le, err := getEnvVal("LOGGING_ENDPOINT", func() (interface{}, error) {
			return os.Stdout, nil
		})
		if err != nil {
			return err
		}
		c.LoggingEndpoint = le.(io.Writer)
		return nil
	}
	redisEndPoint := func(c *Cfg) error {
		re, err := getEnvVal("REDIS_ENDPOINT", func() (interface{}, error) {
			return "localhost:6379", nil
		})
		if err != nil {
			return err
		}
		c.RedisEndpoint = re.(string)
		return nil
	}
	jwtTTL := func(c *Cfg) error {
		re, err := getEnvVal("JWT_TTL", func() (interface{}, error) {
			return (int)(time.Hour.Seconds()), nil
		})
		if err != nil {
			return err
		}
		c.JwtTTL = re.(int)
		return nil
	}
	jwtIss := func(c *Cfg) error {
		re, err := getEnvVal("JWT_ISS", func() (interface{}, error) {
			return "localhost", nil
		})
		if err != nil {
			return err
		}
		c.JwtIss = re.(string)
		return nil
	}
	jwtSub := func(c *Cfg) error {
		re, err := getEnvVal("JWT_SUB", func() (interface{}, error) {
			return "localhost", nil
		})
		if err != nil {
			return err
		}
		c.JwtSub = re.(string)
		return nil
	}
	// var err error
	return NewConfig(privateKey, publicKey, port,
		gitHubClientID, gitHubClientSecret, oauthRedirectURL,
		loggingEndpoint, redisEndPoint,
		jwtTTL, jwtIss, jwtSub)
}

func getEnvVal(key string, defaultValue DefaultValFunc) (interface{}, error) {
	var value interface{}
	var err error
	value = os.Getenv(key)
	if value.(string) == "" {
		log.Printf("Did not set %s - using default", key)
		value, err = defaultValue()
	}
	return value, err
}
