package routes

import (
	"testing"

	"github.com/labstack/echo/v4"
)

func TestLoadRegistersTerraformBackendRoutes(t *testing.T) {
	e := echo.New()
	Load(e)

	expected := map[string]bool{
		"GET /:group/:key":    false,
		"POST /:group/:key":   false,
		"PUT /:group/:key":    false,
		"DELETE /:group/:key": false,
	}

	for _, route := range e.Routes() {
		key := route.Method + " " + route.Path
		if _, ok := expected[key]; ok {
			expected[key] = true
		}
	}

	for route, found := range expected {
		if !found {
			t.Fatalf("expected route %s to be registered", route)
		}
	}
}
