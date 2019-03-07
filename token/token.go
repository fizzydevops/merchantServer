package token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type claims struct {
	identifier string
	stdClaims  jwt.StandardClaims
}

func (c *claims) Valid() error {
	return c.stdClaims.Valid()
}

type token struct {
	key    []byte
	claims *claims
}

func New(identifier string) *token {
	c := &claims{
		identifier: identifier,
		stdClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
			Issuer:    "authServer",
			IssuedAt:  time.Now().Unix(),
		},
	}

	return &token{
		claims: c,
		key:    nil,
	}
}

func (t *token) Key() []byte {
	return t.key
}

func (t *token) SetKey(key []byte) {
	t.key = key
}

func (t *token) Claims() *claims {
	return t.claims
}

func (t *token) SetClaims(c *claims) {
	t.claims = c
}

func (t *token) GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t.claims)
	tokenStr, err := token.SignedString([]byte("webAuth"))
	return tokenStr, err
}

func (t *token) validateToken(tknStr string) (bool, error) {
	tkn, err := jwt.ParseWithClaims(tknStr, &claims{}, func(t *jwt.Token) (interface{}, error){
		return []byte("webAuth"), nil
	})

	if _, ok := tkn.Claims.(*claims); ok && tkn.Valid {
		return true, err
	}

	return false, err
}
