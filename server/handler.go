package server

import (
	"encoding/json"
	"github.com/auth/merchant"
	"github.com/dgrijalva/jwt-go"
	"time"
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
		err := InvalidAuthRequest{missingItems:errMsgs}.Error()
		logMerchantError("Failed to authenticate merchant credentials.", "merchantHandler", err)
		jsonBytes, _:= json.Marshal(map[string]string{
			"status": "error",
			"message": "Insufficient data sent in request.",
			"error": err,
		})
		writer.Write(jsonBytes)
		writer.Flush()
	}

	authenticated, err := merchant.Authenticate(username, password)

	if err != nil {
		logMerchantError("Failed to authenticate merchant credentials.", "merchantHandler", err.Error())
	} else if !authenticated {
		jsonBytes, _ := json.Marshal(map[string]string{
			"status": "Error",
			"message": "Authentication Failure. Invalid credentials.",
		})
		writer.Write(jsonBytes)
		writer.Flush()
	}

	// If authenticated we are going to now get a token for the account.
	// makes a token valid for 60 minutes.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{Issuer:"authServer", ExpiresAt:time.Now().UTC().Add(time.Second*60).Unix()})
	signedStr, err := token.SignedString("test")
}
