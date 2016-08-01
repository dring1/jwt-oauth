package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"net/http"
	"os"

	"github.com/dring1/jwt-oauth/config"
	"github.com/dring1/jwt-oauth/middlewares"
	"github.com/dring1/jwt-oauth/models"
	"github.com/dring1/jwt-oauth/routes"
	"github.com/dring1/jwt-oauth/services"
	"github.com/justinas/alice"
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
			pem.Encode(os.Stdout, &pemBlock)
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
			pem.Encode(os.Stdout, &pemBlock)
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

	gitHubClientId := func(c *config.Cfg) error {
		ghCID, err := getEnvVal("GITHUB_CLIENT_ID", func() (interface{}, error) {
			return nil, errors.Errorf("Did not provide GITHUB_CLIENT_ID")
		})
		if err != nil {
			return err
		}
		c.GitHubClientID = ghCID.(string)
		return nil
	}
	var err error
	c, err = config.NewConfig(privateKey, publicKey, port, gitHubClientId)
	if err != nil {
		log.Fatalf("ERROR: %+v", errors.Wrap(err, "error intializing"))
	}
}

func getEnvVal(key string, defaultValue DefaultValFunc) (interface{}, error) {
	var value interface{}
	var err error
	value = os.Getenv(key)
	if value.(string) == "" {
		value, err = defaultValue()
	}
	return value, err
}

func main() {
	routerC := struct {
		GitHubID string
	}{
		c.GitHubClientID,
	}
	err := tmpl.Execute(w, routerC)
	router := routes.NewRouter()
	chain := alice.New(middlewares.LoggingHandler, middlewares.RecoverHandler).Then(router)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), chain)
	log.Fatal(err)
}
