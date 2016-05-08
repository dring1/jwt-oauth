package authentication

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dring1/orm/config"
	"github.com/dring1/orm/lib/cache"
	"github.com/dring1/orm/models"
)

type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	Cache      *Cache.Cache
}

const (
	TokenDuration = 72
	ExpireOffset  = 3600
)

var JWTBackendInstance *JWTAuthenticationBackend

func init() {
	abi, err := NewJWTBackend()
	if err != nil {
		log.Fatal(err)
	}
	JWTBackendInstance = abi
}

func NewJWTBackend() (*JWTAuthenticationBackend, error) {
	rawPrivData, err := ioutil.ReadFile(config.Cfg.PrivateKeyPath)
	if err != nil {
		return nil, err
	}
	log.Println("Here")
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
	ab := &JWTAuthenticationBackend{
		privateKey: privateKey,
		PublicKey:  publicKey,
		Cache:      Cache.NewCache(),
	}
	return ab, nil
}

func (backend *JWTAuthenticationBackend) GenerateToken(userUUID string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims["exp"] = time.Now().Add(time.Hour * time.Duration(config.Cfg.JWTExpirationDelta)).Unix()
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["sub"] = userUUID
	log.Println(backend)
	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

// Actually auth
func (backend *JWTAuthenticationBackend) Authenticate(user *models.User) bool {
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
