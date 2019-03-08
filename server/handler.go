package server

import (
	"encoding/base64"
	"github.com/auth/merchant"
	"strings"
)

// merchantHandler will either validate credentials and send back a token or just validate a token.
func merchantHandler(data map[string]interface{}) {
	reqType, ok := data["type"].(string)
	var err error

	if !ok {
		err = &InsufficientDataError{[]string{"No type provided in request."}}
		logServerError("Invalid request.", "merchantHandler", err.Error())
		writeResponse(map[string]interface{}{"status": "error", "message": "Insufficient data sent in request", "error": err.Error()})
		return
	} else if strings.ToLower(reqType) != "auth" || strings.ToLower(reqType) != "validate" {
		err = &InvalidRequestTypeError{reqType: reqType}
		logServerError("Invalid request.", "merchantHandler", err.Error())
		writeResponse(map[string]interface{}{"status": "error", "message": "Invalid type.", "error": err.Error()})
		return
	}

	if strings.ToLower(reqType) == "auth" {
		authenticateMerchant(data)
	} else {
		validateToken(data)
	}
}

// validateMerchant validates if the incoming request token is valid for use.
func validateToken(data map[string]interface{}) {
	var errMsgs []string

	username, ok := data["username"].(string)

	if !ok {
		errMsgs = append(errMsgs, "No username provided in request.")
	} else if username == "" {
		errMsgs = append(errMsgs,"Username cannot be zero value.")
	}

	token, ok := data["token"].(string)

	if !ok {
		errMsgs = append(errMsgs, "No token provided in request.")
	} else if token == "" {
		errMsgs = append(errMsgs, "Token cannot be zero value.")
	}

	if errMsgs != nil {
		err := &InsufficientDataError{missingItems: errMsgs}
		logServerError("Failed to authenticate merchant credentials.", "validateToken", err.Error())
		writeResponse(map[string]interface{}{"status": "error", "message": "Insufficient data sent in request.", "error": err.Error()})
		return
	}

	valid, err := ValidateToken(token, username)

	if err != nil {
		logServerError("Failed to validate token.", "validateToken", err.Error())
		writeResponse(map[string]interface{}{"status": "error", "message": "Failed to authenticate token.", "error": err.Error()})
	} else if !valid {
		writeResponse(map[string]interface{}{"status": "error", "message": "Authentication Failure. Token expired please request a new one."})
	}

	writeResponse(map[string]interface{}{"status": "success", "message": "Successfully authenticated."})
}

// authenticateMerchant handles with authenticating user credentials and if successfully authenticated, grant a jwt token.
func authenticateMerchant(data map[string]interface{}) {
	// Validate incoming request
	var errMsgs []string

	username, ok := data["username"].(string)

	if !ok {
		errMsgs = append(errMsgs, "No username provided in request.")
	} else if username == "" {
		errMsgs = append(errMsgs, "Username cannot be zero value.")
	}

	password, ok := data["password"].(string)

	if !ok {
		errMsgs = append(errMsgs, "No password provided in request.")
	} else if password == "" {
		errMsgs = append(errMsgs, "Password cannot be zero value.")
	}

	passwordBytes, err := base64.StdEncoding.DecodeString(password)

	if err != nil {
		errMsgs = append(errMsgs, err.Error())
	}

	if errMsgs != nil {
		err := &InsufficientDataError{missingItems: errMsgs}
		logServerError("Failed to authenticate merchant credentials.", "authenticateMerchant", err.Error())
		writeResponse(map[string]interface{}{"status": "error", "message": "Insufficient data sent in request.", "error": err.Error()})
		return
	}

	m := merchant.New(username, passwordBytes)

	authenticated, err := m.Authenticate()

	if err != nil {
		logServerError("Failed to authenticate merchant credentials.", "authenticateMerchant", err.Error())
		writeResponse(map[string]interface{}{"status": "error", "message": "Failed to authenticate merchant credentials."})
	} else if !authenticated {
		writeResponse(map[string]interface{}{"status": "error", "message": "Authentication Failure. Invalid credentials."})
		return
	}

	// If authenticated we are going to now get a token for the account.
	// makes a token valid for 60 minutes.
	token, err := GenerateToken(username)

	if err != nil {
		logServerError("Failed to generate a token.", "merchantHandler", err.Error())
	}

	writeResponse(map[string]interface{}{
		"status":   "success",
		"message":  "Successfully authenticated.",
		"username": username,
		"token":    token,
	})
}
