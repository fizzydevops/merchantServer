package user

import (
	"auth/db"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	username string
	password []byte
	token    string
}

func New(username string) *user {
	return &user{
		username: username,
	}
}

func (u *user) Username() string {
	return u.username
}

func (u *user) SetMerchant(username string) {
	u.username = username
}

func (u *user) Password() []byte {
	return u.password
}

func (u *user) SetPassword(password []byte) {
	u.password = password
}

func (u *user) Token() string {
	return u.token
}

func (u *user) setToken(token string) {
	u.token = token
}

// Authenticate connect to the database to authenticate user credentials
func Authenticate(username string, password []byte) (bool, error) {
	conn := db.New("merchantdb")

	results, err := conn.QueryAndScan(`SELECT password FROM login WHERE username = ? `, []interface{}{username})

	if err != nil {
		return false, err
	}

	if len(results) == 0 {
		err = &UsernameNotFoundError{username:username}
		return false, err
	}

	err = bcrypt.CompareHashAndPassword(results["password"].([]byte), password)

	// If err is returned from bcrypt.CompareHashAndPassword it means passwords did not match.
	if err != nil {
		return false, nil
	}

	return true, nil
}
