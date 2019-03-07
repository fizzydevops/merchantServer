package merchant

import (
	"github.com/Auth/db"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type merchant struct {
	username string
	password []byte
	token    string
}

func New(username string, password []byte) *merchant {
	return &merchant{
		username: username,
		password: password,
	}
}

func (m *merchant) Username() string {
	return m.username
}

func (m *merchant) SetMerchant(username string) {
	m.username = username
}

func (m *merchant) Password() []byte {
	return m.password
}

func (m *merchant) SetPassword(password []byte) {
	m.password = password
}

func (m *merchant) Token() string {
	return m.token
}

func (m *merchant) setToken(token string) {
	m.token = token
}

// Authenticate connect to the database to authenticate merchant credentials
func (m *merchant) Authenticate(username string, password []byte) (bool, error) {
	conn, err := db.NewConnection("merchantdb")

	if err != nil {
		log.Println(map[string]string{
			"status":   "error",
			"message":  "Failed to establish database connection.",
			"database": "merchantdb",
			"function": "Authenticate",
			"package":  "merchant",
			"error":    err.Error(),
		})
		return false, err
	}

	results, err := conn.QueryAndScan(`SELECT password FROM login WHERE username = ? `, []interface{}{username})

	if err != nil {
		log.Println(map[string]string{
			"status":   "error",
			"message":  "Failed to query database.",
			"database": "merchantdb",
			"function": "Authenticate",
			"package":  "merchant",
			"error":    err.Error(),
		})
	}

	err = bcrypt.CompareHashAndPassword(results["password"].([]byte), password)

	// If err is returned from bcrypt it means passwords did not match.
	if err != nil {
		return false, nil
	}

	return true, nil
}
