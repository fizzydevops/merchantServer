package server

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GenerateToken(identifier string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		Issuer:    "auth",
		IssuedAt:  time.Now().Unix(),
	})
	tokenStr, err := token.SignedString([]byte(identifier))
	return tokenStr, err
}

func ValidateToken(tknStr, identifier string) (bool, error) {
	tkn, err := jwt.ParseWithClaims(tknStr, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(identifier), nil
	})

	if _, ok := tkn.Claims.(*jwt.StandardClaims); ok && tkn.Valid {
		return true, err
	}

	return false, err
}
