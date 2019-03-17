package server

import (
	"bufio"
	"encoding/base64"
	"auth/user"
	"strings"
)

//userHandler will either validate credentials and send back a token or just validate a token.
func userHandler(handlerStruct *handler) {
	var err error
	data := handlerStruct.Data
	writer := handlerStruct.Writer

	reqType, ok := data["type"].(string)

	if !ok {
		err = &InsufficientDataError{[]string{"No type provided in request."}}
		logger.Log(map[string]interface{}{
			"file":     "handler.go",
			"package":  "server",
			"function": "userHandler",
			"message":  "Invalid request.",
			"error":    err.Error(),
		})
		writeResponse(writer, map[string]interface{}{"status": "error", "message": "Insufficient data sent in request", "error": err.Error()})
		return
	} else if strings.ToLower(reqType) != "auth" && strings.ToLower(reqType) != "validate" {
		err = &InvalidRequestTypeError{reqType: reqType}
		logger.Log(map[string]interface{}{
			"file":     "handler.go",
			"package":  "server",
			"function": "userHandler",
			"message":  "Invalid request.",
			"error":    err.Error(),
		})
		writeResponse(writer, map[string]interface{}{"status": "error", "message": "Invalid type.", "error": err.Error()})
		return
	}

	if strings.ToLower(reqType) == "auth" {
		authenticateUser(writer, data)
	} else {
		validateToken(writer, data)
	}
}

// validateToken validates if the incoming request token is valid for use.
func validateToken(writer *bufio.Writer, data map[string]interface{}) {
	var errMsgs []string

	username, ok := data["username"].(string)

	if !ok {
		errMsgs = append(errMsgs, "No username provided in request.")
	} else if username == "" {
		errMsgs = append(errMsgs, "Username cannot be zero value.")
	}

	token, ok := data["token"].(string)

	if !ok {
		errMsgs = append(errMsgs, "No token provided in request.")
	} else if token == "" {
		errMsgs = append(errMsgs, "Token cannot be zero value.")
	}

	if errMsgs != nil {
		err := &InsufficientDataError{missingItems: errMsgs}
		logger.Log(map[string]interface{}{
			"file":     "handler.go",
			"package":  "server",
			"function": "validateToken",
			"message":  "Failed to authenticate user credentials.",
			"error":    err.Error(),
		})
		writeResponse(writer, map[string]interface{}{"status": "error", "message": "Insufficient data sent in request.", "error": err.Error()})
		return
	}

	valid, err := ValidateToken(token, username)

	if err != nil {
		logger.Log(map[string]interface{}{
			"file":     "handler.go",
			"package":  "server",
			"function": "validateToken",
			"message":  "Failed to validate token.",
			"error":    err.Error(),
		})
		writeResponse(writer, map[string]interface{}{"status": "error", "message": "Failed to authenticate token.", "error": err.Error()})
		return
	} else if !valid {
		writeResponse(writer, map[string]interface{}{"status": "error", "message": "Authentication Failure. Token expired please request a new one."})
		return
	}

	writeResponse(writer, map[string]interface{}{"status": "success", "message": "Successfully validated token."})
}

// authenticateUser handles with authenticating user credentials and if successfully authenticated, grant a jwt token.
func authenticateUser(writer *bufio.Writer, data map[string]interface{}) {
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
		logger.Log(map[string]interface{}{
			"file":     "handler.go",
			"package":  "server",
			"function": "authenticateUser",
			"message":  "Failed to validate token.",
			"error":    err.Error(),
		})
		writeResponse(writer, map[string]interface{}{"status": "error", "message": "Insufficient data sent in request.", "error": err.Error()})
		return
	}

	authenticated, err := user.Authenticate(username, passwordBytes)

	if err != nil {
		logger.Log(map[string]interface{}{
			"file":     "handler.go",
			"package":  "server",
			"function": "authenticateUser",
			"message":  "Failed to authenticate user credentials.",
			"error":    err.Error(),
		})
		writeResponse(writer, map[string]interface{}{"status": "error", "message": "Failed to authenticate user credentials."})
		return
	} else if !authenticated {
		writeResponse(writer, map[string]interface{}{"status": "error", "message": "Authentication Failure. Invalid credentials."})
		return
	}

	// If authenticated we are going to now get a token for the account.
	// makes a token valid for 60 minutes.
	token, err := GenerateToken(username)

	if err != nil {
		logger.Log(map[string]interface{}{
			"file":     "handler.go",
			"package":  "server",
			"function": "authenticateUser",
			"message":  "Failed to generate a token.",
			"error":    err.Error(),
		})
		return
	}

	writeResponse(writer, map[string]interface{}{
		"status":   "success",
		"message":  "Successfully authenticated.",
		"username": username,
		"token":    token,
	})
}
