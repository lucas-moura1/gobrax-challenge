package handler

import (
	"encoding/json"
	"net/http"
)

func errorHandler(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
