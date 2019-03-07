package server

import (
	"encoding/json"
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
		logMerchantError("Failed to authenticate credentials.", "merchantHandler", err)
		jsonBytes, _:= json.Marshal(map[string]string{
			"status": "error",
			"message": "Insufficient data sent in request.",
			"error": err,
		})
		writer.Write(jsonBytes)
		writer.Flush()
	}



}
