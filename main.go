package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/{group}/{key}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		group := r.PathValue("group")
		key := r.PathValue("key")
		path := getStatePath(group, key)

		if !isStateExists(path) && r.Method != "POST" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{\"error\":\"Status not found\"}"))
			return
		}

		switch r.Method {
		case "GET":
			content, err := os.ReadFile(path)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("{\"error\":\"" + err.Error() + "\"}"))
				return
			}
			w.Write(content)
		case "POST":
			defer r.Body.Close()
			content, _ := io.ReadAll(r.Body)

			if !json.Valid(content) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("{\"error\":\"Invalid JSON format\"}"))
				return
			}

			err := os.WriteFile(path, content, 0666)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("{\"error\":\"" + err.Error() + "\"}"))
				return
			}
			w.Write(content)
		case "PUT":
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"error\":\"Lock is not supported\"}"))
		case "DELETE":
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"error\":\"Unlock is not supported\"}"))
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"error\":\"Method is not supported\"}"))
		}
	})

	http.ListenAndServe(":8080", nil)
}

func getStatePath(group string, key string) string {
	return fmt.Sprintf("storage/%s-%s.json", group, key)

}

func isStateExists(path string) bool {
	_, error := os.Stat(path)
	return !errors.Is(error, os.ErrNotExist)
}
