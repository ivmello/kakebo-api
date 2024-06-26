package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data  any `json:"data"`
	Error any `json:"error"`
}

func JSONResponse(w http.ResponseWriter, output any, status int) {
	jsonResponse := Response{
		Data: output,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func JSONErrorResponse(w http.ResponseWriter, message string, status int) {
	jsonResponse := Response{
		Error: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		http.Error(w, err.Error(), status)
	}
}
