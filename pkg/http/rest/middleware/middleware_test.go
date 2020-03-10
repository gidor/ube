package middleware

import (
	"reflect"
	"testing"

	"github.com/gidor/ube/test"
)

func TestCreateMiddleware(t *testing.T) {
	logger, _ := test.LoggerMock()
	mw := CreateMiddleware(logger)

	expectedType := reflect.TypeOf(&Middleware{})
	if r := reflect.TypeOf(mw); r != expectedType {
		t.Errorf("middleware has wrong type: got %v want %v",
			r, expectedType)
	}
}
