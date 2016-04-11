package authentication

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dring1/orm/config"
	"github.com/dring1/orm/models"
)

type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

const (
	TokenDuration = 72
	ExpireOffset  = 3600
)

var authBackendInstance *JWTAuthenticationBackend

func JWTBackend() (*JWTAuthenticationBackend, error) {
	rawPrivData, err := ioutil.ReadFile(config.Cfg.PrivateKeyPath)
	if err != nil {
		return nil, err
	}
	privateKey, err := getPrivateKey(rawPrivData)
	if err != nil {
		return nil, err
	}
	rawPubData, err := ioutil.ReadFile(config.Cfg.PublicKeyPath)
	if err != nil {
		return nil, err
	}
	publicKey, err := getPublicKey(rawPubData)
	if err != nil {
		return nil, err
	}
	authBackendInstance = &JWTAuthenticationBackend{
		privateKey: privateKey,
		PublicKey:  publicKey,
	}
	return authBackendInstance, nil
}

func (backend *JWTAuthenticationBackend) GenerateToken(userUUID string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims["exp"] = time.Now().Add(time.Hour * time.Duration(config.Cfg.JWTExpirationDelta)).Unix()
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["sub"] = userUUID
	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func (backend *JWTAuthenticationBackend) Authenticate(user *models.User) bool {
	return true
}

func (backend *JWTAuthenticationBackend) TimeToExpire(timestamp interface{}) int64 {

	if ts, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(ts), 0)
		if remainder := tm.Sub(time.Now()); remainder > 0 {
			return int64(remainder.Seconds()) + ExpireOffset
		}
	}
	return ExpireOffset
}

func getPrivateKey(rawPemData []byte) (*rsa.PrivateKey, error) {
	data, _ := pem.Decode([]byte(rawPemData))
	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKeyImported, nil
}

func getPublicKey(rawPemData []byte) (*rsa.PublicKey, error) {
	data, _ := pem.Decode([]byte(rawPemData))
	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)
	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		return nil, fmt.Errorf("Not a valid RSA public key")
	}

	return rsaPub, nil
}
