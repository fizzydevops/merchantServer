package logger_test

import (
	"github.com/auth/logger"
	"testing"
)

func TestLog(t *testing.T) {
	data := map[string]interface{}{
		"status":   "error",
		"message":  "test message",
		"function": "TestLog",
		"package":  "logger_test",
		"error":    "This is a test.",
	}

	logger.Log(data)
}
