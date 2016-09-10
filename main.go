package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/dring1/jwt-oauth/cache"
	"github.com/dring1/jwt-oauth/config"
	"github.com/dring1/jwt-oauth/controllers"
	"github.com/dring1/jwt-oauth/database"
	"github.com/dring1/jwt-oauth/middleware"
	"github.com/dring1/jwt-oauth/routes"
	"github.com/pkg/errors"
)

var c *config.Cfg

type DefaultValFunc func() (interface{}, error)

func init() {
	var PrivateKey *pem.Block
	privateKey := func(c *config.Cfg) error {
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
	publicKey := func(c *config.Cfg) error {
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

	port := func(c *config.Cfg) error {
		p, err := getEnvVal("PORT", func() (interface{}, error) {
			return 8080, nil
		})
		if err != nil {
			return err
		}
		c.Port = p.(int)
		return nil
	}

	gitHubClientID := func(c *config.Cfg) error {
		ghCID, err := getEnvVal("GITHUB_CLIENT_ID", func() (interface{}, error) {
			return nil, errors.Errorf("Did not provide GITHUB_CLIENT_ID")
		})
		if err != nil {
			return err
		}
		c.GitHubClientID = ghCID.(string)
		return nil
	}

	gitHubClientSecret := func(c *config.Cfg) error {
		ghCS, err := getEnvVal("GITHUB_CLIENT_SECRET", func() (interface{}, error) {
			return nil, errors.Errorf("Did not provide GITHUB_CLIENT_SECRET")
		})
		if err != nil {
			return err
		}
		c.GitHubClientSecret = ghCS.(string)
		return nil
	}

	oauthRedirectURL := func(c *config.Cfg) error {
		rdURL, err := getEnvVal("OAUTH_REDIRECT_URL", func() (interface{}, error) {
			return fmt.Sprintf("http://localhost:%d/github/callback", c.Port), nil
		})
		if err != nil {
			return err
		}
		c.OauthRedirectURL = rdURL.(string)
		return nil
	}
	loggingEndpoint := func(c *config.Cfg) error {
		le, err := getEnvVal("LOGGING_ENDPOINT", func() (interface{}, error) {
			return os.Stdout, nil
		})
		if err != nil {
			return err
		}
		c.LoggingEndpoint = le.(io.Writer)
		return nil
	}
	redisEndPoint := func(c *config.Cfg) error {
		re, err := getEnvVal("REDIS_ENDPOINT", func() (interface{}, error) {
			return "localhost:6379", nil
		})
		if err != nil {
			return err
		}
		c.RedisEndpoint = re.(string)
		return nil
	}
	var err error
	c, err = config.NewConfig(privateKey, publicKey, port,
		gitHubClientID, gitHubClientSecret, oauthRedirectURL, loggingEndpoint, redisEndPoint)
	if err != nil {
		log.Fatalf("ERROR: %+v", errors.Wrap(err, "error intializing"))
	}
	log.Println("Done initializing config")
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

func main() {
	// Init services
	db, _ := database.NewDatabaseService()
	ch, _ := cache.New(c.RedisEndpoint)

	// Init controllers
	ctrls := controllers.New(db, ch)

	// Init router
	router := routes.New(c.GitHubClientID, c.GitHubClientSecret, c.OauthRedirectURL, ctrls)

	// Apply middlewares
	middlewares := []middleware.Middleware{
		middleware.NewApacheLoggingHandler(c.LoggingEndpoint),
	}
	middlewares = append(middlewares, middleware.DefaultMiddleWare()...)

	log.Printf("Serving on port :%d", c.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", c.Port), middleware.Handlers(router, middlewares...))
	log.Fatal(err)
}
