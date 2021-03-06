package rest

import (
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

func TestCreateHandler(t *testing.T) {
	// logger, _ := test.LoggerMock()
	handler := CreateHandler()

	expectedRouterType := reflect.TypeOf(mux.NewRouter())
	if r := reflect.TypeOf(handler.GetRouter()); r != expectedRouterType {
		t.Errorf("handler has wrong type of router: got %v want %v",
			r, expectedRouterType)
	}
}
