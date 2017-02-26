package models

import _jwt "github.com/dgrijalva/jwt-go"

type Token struct {
	T           _jwt.Token `json:"-"`
	TokenString string     `json:"token"`
}
