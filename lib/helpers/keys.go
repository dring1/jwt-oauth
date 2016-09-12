package helpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func GenPrivateKey() *pem.Block {

	pk, _ := rsa.GenerateKey(rand.Reader, 1024)
	bits := x509.MarshalPKCS1PrivateKey(pk)
	pemBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: bits,
	}
	return &pemBlock
}

func GenPublicKey(pKey *pem.Block) (*pem.Block, error) {

	privKeyBytes := pKey.Bytes
	privKey, err := x509.ParsePKCS1PrivateKey(privKeyBytes)
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
}
