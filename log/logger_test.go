package log_test

import (
	"github.com/auth/log"
	"testing"
)

func TestLog(t *testing.T) {

	for i := 0; i < 1000; i++ {
		data := map[string]interface{}{
			"status":   "error",
			"message":  "test message",
			"function": "TestLog",
			"package":  "logger_test",
			"error":    "This is a test.",
		}
		log.Log(data)
		t.Log(i)
	}
}
