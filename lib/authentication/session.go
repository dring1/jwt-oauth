package authentication

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dring1/orm/config"
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

func InitJWTAuthenticationBackend() (*JWTAuthenticationBackend, error) {
	if authBackendInstance == nil {
		privateKey, err := getPrivateKey()
		if err != nil {
			return nil, err
		}
		publicKey, err := getPublicKey()
		if err != nil {
			return nil, err
		}
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: privateKey,
			PublicKey:  publicKey,
		}
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

func getPrivateKey() (*rsa.PrivateKey, error) {
	privateKeyFile, err := os.Open(config.Cfg.PrivateKeyPath)
	defer privateKeyFile.Close()
	if err != nil {
		return nil, err
	}
	pemfileInfo, _ := privateKeyFile.Stat()
	size := pemfileInfo.Size()
	pemBytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pemBytes)
	if err != nil {
		return nil, err
	}

	data, _ := pem.Decode(pemBytes)
	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKeyImported, nil
}

func getPublicKey() (*rsa.PublicKey, error) {
	publicKeyFile, err := os.Open(config.Cfg.PublicKeyPath)
	defer publicKeyFile.Close()
	if err != nil {
		return nil, err
	}

	pemfileinfo, _ := publicKeyFile.Stat()
	size := pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)
	if err != nil {
		return nil, err
	}

	data, _ := pem.Decode([]byte(pembytes))

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
