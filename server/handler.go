package server

import (
	"encoding/base64"
	"github.com/auth/merchant"
	"log"
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
		err = &InvalidRequestTypeError{reqType:reqType}
		logServerError("Invalid request.", "merchantHandler", err.Error())
		writeResponse(map[string]interface{}{"status": "error", "message": "Invalid type.", "error": err.Error()})
		return
	}

	if strings.ToLower(reqType) == "auth" {
		authenticateMerchant(data)
	} else {
		validateMerchant(data)
	}
}

// validateMerchant validates if the incoming request token is valid for use.
func validateMerchant(data map[string]interface{}) {

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

	password, err := base64.StdEncoding.DecodeString(data["password"].(string))

	if err != nil {
		log.Println(err.Error())
	}

	if !ok {
		errMsgs = append(errMsgs, "No password provided in request.")
	} else if password == nil {
		errMsgs = append(errMsgs, "Password cannot be nil.")
	}

	if errMsgs != nil {
		err := &InsufficientDataError{missingItems: errMsgs}
		logServerError("Failed to authenticate merchant credentials.", "merchantHandler", err.Error())
		writeResponse(map[string]interface{}{"status":  "error", "message": "Insufficient data sent in request.", "error":  err.Error()})
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
	token, err := GenerateToken(username)

	if err != nil {
		logServerError("Failed to generate a token", "merchantHandler", err.Error())
	}

	writeResponse(map[string]interface{}{
		"status":   "success",
		"message":  "Successfully authenticated.",
		"username": username,
		"token":    token,
	})
}
