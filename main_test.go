package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPortUsesDefaultPort(t *testing.T) {
	t.Setenv("HTTP_PORT", "")

	if got := getPort(); got != ":8080" {
		t.Fatalf("expected default port :8080, got %q", got)
	}
}

func TestGetPortUsesConfiguredPort(t *testing.T) {
	t.Setenv("HTTP_PORT", "9090")

	if got := getPort(); got != ":9090" {
		t.Fatalf("expected configured port :9090, got %q", got)
	}
}

func TestUseBasicAuthIsDisabledWhenCredentialsAreUnset(t *testing.T) {
	t.Setenv("BASIC_AUTH_USERNAME", "")
	t.Setenv("BASIC_AUTH_PASSWORD", "")

	response := serveWithBasicAuth("", "")
	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestUseBasicAuthRejectsMissingCredentials(t *testing.T) {
	t.Setenv("BASIC_AUTH_USERNAME", "user")
	t.Setenv("BASIC_AUTH_PASSWORD", "pass")

	response := serveWithBasicAuth("", "")
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}
}

func TestUseBasicAuthAcceptsConfiguredCredentials(t *testing.T) {
	t.Setenv("BASIC_AUTH_USERNAME", "user")
	t.Setenv("BASIC_AUTH_PASSWORD", "pass")

	response := serveWithBasicAuth("user", "pass")
	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func serveWithBasicAuth(username string, password string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := useBasicAuth(mux)

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	if username != "" || password != "" {
		request.SetBasicAuth(username, password)
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	return response
}
