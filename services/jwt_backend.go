package services

import (
	"encoding/pem"
	"log"
	"os"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type JWTAuthenticationBackend struct {
	privateKey, PublicKey []byte
}

const (
	TokenDuration = 72
	ExpireOffset  = 3600
)

var (
	jwtOnce            sync.Once
	jwtBackendInstance *JWTAuthenticationBackend
)

func JWTBackend() *JWTAuthenticationBackend {
	jwtOnce.Do(func() {

		abi, err := NewJWTBackend()
		if err != nil {
			log.Fatal(err)
		}
		jwtBackendInstance = abi

	})
	return jwtBackendInstance
}

func NewJWTBackend() (*JWTAuthenticationBackend, error) {
	rawPrivData := []byte(os.Getenv("PRIVATE_KEY")) // ioutil.ReadFile(config.Cfg.PrivateKeyPath)
	// if err != nil {
	// 	return nil, err
	// }
	privateKey, err := getPrivateKey(rawPrivData)
	if err != nil {
		return nil, err
	}
	rawPubData := []byte(os.Getenv("PUBLIC_KEY")) //ioutil.ReadFile(config.Cfg.PublicKeyPath)
	// if err != nil {
	// 	return nil, err
	// }

	publicKey, err := getPublicKey(rawPubData)
	if err != nil {
		return nil, err
	}
	ab := &JWTAuthenticationBackend{
		privateKey: privateKey,
		PublicKey:  publicKey,
	}
	return ab, nil
}

type CustomClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func (backend *JWTAuthenticationBackend) GenerateToken(userID string) (string, error) {
	exp := time.Now().Add(time.Hour * 1).Unix()
	iss := "jwt-oauth.com"
	sub := "jwt-oauth"
	claims := CustomClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: exp,
			Issuer:    iss,
			Subject:   sub,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Insert into cache here ?
func (backend *JWTAuthenticationBackend) Authenticate(interface{}) bool {
	return true
}

func (backend *JWTAuthenticationBackend) Logout(tokenString string, token *jwt.Token) error {
	return nil
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

func getPrivateKey(rawPemData []byte) ([]byte, error) {
	data, _ := pem.Decode([]byte(rawPemData))
	return data.Bytes, nil
}

func getPublicKey(rawPemData []byte) ([]byte, error) {
	data, _ := pem.Decode([]byte(rawPemData))
	return data.Bytes, nil
}
