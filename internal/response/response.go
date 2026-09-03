package response

import (
	"encoding/json"
	"log"
	"net/http"
)

// WithError sends a JSON error response with the given HTTP status code.
func WithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	WithJSON(w, code, errResponse{Error: msg})
}

// WithJSON converts any Go value to JSON and writes it to the HTTP response.
func WithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Failed to Marshal JSON response: %v", payload)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
