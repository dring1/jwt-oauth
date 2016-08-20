package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dring1/jwt-oauth/config"
	"github.com/dring1/jwt-oauth/middleware"
	"github.com/dring1/jwt-oauth/models"
	"github.com/dring1/jwt-oauth/routes"
	"github.com/dring1/jwt-oauth/services"
	"github.com/pkg/errors"
)

var c *config.Cfg

type DefaultValFunc func() (interface{}, error)

func init() {
	services.Database()
	services.Database().HasTable(&models.User{})
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
	var err error
	c, err = config.NewConfig(privateKey, publicKey, port,
		gitHubClientID, gitHubClientSecret, oauthRedirectURL)
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
	router := routes.New(c.GitHubClientID, c.GitHubClientSecret)
	middlewares := []middleware.Middleware{
		middleware.NewApacheLoggingHandler(os.Stdout),
	}
	middlewares = append(middlewares, middleware.DefaultMiddleWare()...)

	log.Printf("Serving on port :%d", c.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", c.Port), middleware.Handlers(router, middlewares...))
	log.Fatal(err)
}
