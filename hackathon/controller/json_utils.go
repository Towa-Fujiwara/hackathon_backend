package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

// respondJSON はJSON形式でレスポンスを返します。
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("failed to encode json: %v", err)
		}
	}
}

// decodeBody はリクエストボディをデコードします。
func decodeBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}