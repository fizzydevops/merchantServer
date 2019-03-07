package server

import (
	"github.com/Auth/token"
	"github.com/auth/merchant"
)

func merchantHandler(data map[string]interface{}) {
	// Validate incoming request
	var errMsgs []string

	username, ok := data["username"].(string)

	if !ok {
		errMsgs = append(errMsgs, "No username provided in request.")
	} else if username == "" {
		errMsgs = append(errMsgs, "Username cannot be zero value.")
	}

	password, ok := data["password"].([]byte)

	if !ok {
		errMsgs = append(errMsgs, "No password provided in request.")
	} else if password == nil {
		errMsgs = append(errMsgs, "Password cannot be nil.")
	}

	if errMsgs != nil {
		err := InvalidAuthRequest{missingItems: errMsgs}.Error()
		logServerError("Failed to authenticate merchant credentials.", "merchantHandler", err)
		writeResponse(map[string]interface{}{"status":  "error", "message": "Insufficient data sent in request.", "error":  err})
		return
	}

	m := merchant.New(username, password)
	authenticated, err := m.Authenticate()

	if err != nil {
		logServerError("Failed to authenticate merchant credentials.", "merchantHandler", err.Error())
		writeResponse(map[string]interface{}{"status": "error", "message": "Failed to authenticate merchant credentials."})
	} else if !authenticated {
		writeResponse(map[string]interface{}{"status": "error", "message": "Authentication Failure. Invalid credentials."})
		return
	}

	// If authenticated we are going to now get a token for the account.
	// makes a token valid for 60 minutes.
	t := token.New(username)
	tokenStr, err := t.GenerateToken()

	if err != nil {
		logServerError("Failed to generate a token", "merchantHandler", err.Error())
	}

	writeResponse(map[string]interface{}{
		"status":   "success",
		"message":  "Successfully authenticated.",
		"username": username,
		"token":    tokenStr,
	})
}
