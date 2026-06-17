package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"terraform-state-http-backend/routes"
	"time"
)

func main() {
	mux := http.NewServeMux()
	routes.Load(mux)

	handler := useBasicAuth(mux)
	handler = logRequests(handler)

	log.Fatal(http.ListenAndServe(getPort(), handler))
}

func getPort() string {
	port := os.Getenv("HTTP_PORT")

	if port == "" {
		return ":8080"
	}

	return fmt.Sprintf(":%s", port)
}

func useBasicAuth(next http.Handler) http.Handler {
	username := os.Getenv("BASIC_AUTH_USERNAME")
	password := os.Getenv("BASIC_AUTH_PASSWORD")
	if username == "" || password == "" {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inputUsername, inputPassword, ok := r.BasicAuth()
		if !ok || inputUsername != username || inputPassword != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="terraform-state-http-backend"`)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s, (%s)", r.Method, r.URL.RequestURI(), time.Since(start))
	})
}
