// json.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

// respondWithJSON — converts any Go value to JSON and writes it to the response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	data, err := json.Marshal(payload) // Go struct → JSON bytes

	if err != nil {
		w.WriteHeader(500)
		log.Printf("Failed to Marshal JSON response : %v", payload)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code) // set HTTP status code
	w.Write(data)       // write the body

}
