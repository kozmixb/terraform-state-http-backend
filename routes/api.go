package routes

import (
	stdhttp "net/http"
	"terraform-state-http-backend/http"
)

func Load(mux *stdhttp.ServeMux) {
	mux.HandleFunc("GET /{group}/{key}", http.Show)
	mux.HandleFunc("POST /{group}/{key}", http.Update)
	mux.HandleFunc("PUT /{group}/{key}", http.Lock)
	mux.HandleFunc("DELETE /{group}/{key}", http.Unlock)
}
