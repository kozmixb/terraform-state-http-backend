package http

import (
	"encoding/json"
	"errors"
	"io"
	stdhttp "net/http"
	"terraform-state-http-backend/app"
	"terraform-state-http-backend/app/drivers"
)

func Show(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	result, err := app.Show(r.PathValue("group"), r.PathValue("key"))

	if errors.Is(err, drivers.ErrNotFound) {
		writeJSON(w, stdhttp.StatusNotFound, map[string]interface{}{
			"error": "Not Found",
		})
		return
	}
	if err != nil {
		writeJSON(w, stdhttp.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	writeBlob(w, stdhttp.StatusOK, result)
}

func Update(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var body []byte
	if r.Body != nil {
		var err error
		body, err = io.ReadAll(r.Body)
		if err != nil {
			writeJSON(w, stdhttp.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
	}

	result, err := app.Update(r.PathValue("group"), r.PathValue("key"), body)
	if err != nil {
		writeJSON(w, stdhttp.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	writeBlob(w, stdhttp.StatusOK, result)
}

func Lock(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, stdhttp.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	result, err := app.Lock(r.PathValue("group"), r.PathValue("key"), body)
	if errors.Is(err, drivers.ErrLocked) {
		writeBlob(w, stdhttp.StatusLocked, result)
		return
	}
	if err != nil {
		writeJSON(w, stdhttp.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	writeBlob(w, stdhttp.StatusOK, result)
}

func Unlock(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, stdhttp.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	err = app.Unlock(r.PathValue("group"), r.PathValue("key"), body)
	if errors.Is(err, drivers.ErrUnlockMismatch) {
		writeJSON(w, stdhttp.StatusConflict, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	if err != nil {
		writeJSON(w, stdhttp.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(stdhttp.StatusOK)
}

func writeBlob(w stdhttp.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func writeJSON(w stdhttp.ResponseWriter, status int, body map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
