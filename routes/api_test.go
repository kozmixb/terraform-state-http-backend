package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoadRegistersTerraformBackendRoutes(t *testing.T) {
	mux := http.NewServeMux()
	Load(mux)

	tests := []struct {
		method string
		status int
	}{
		{method: http.MethodGet, status: http.StatusNotFound},
		{method: http.MethodPost, status: http.StatusOK},
		{method: http.MethodPut, status: http.StatusOK},
		{method: http.MethodDelete, status: http.StatusOK},
	}

	t.Setenv("DRIVER", "file")
	t.Chdir(t.TempDir())

	for _, test := range tests {
		request := httptest.NewRequest(test.method, "/group/key", nil)
		response := httptest.NewRecorder()

		mux.ServeHTTP(response, request)

		if response.Code != test.status {
			t.Fatalf("expected %s status %d, got %d", test.method, test.status, response.Code)
		}
	}
}
