package auth

import (
	"errors"
	"net/http"
	"strings"
)

// This function extracts and APIKey from the headers of an HTTP request

// Example : apikey <key>
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}

	vals := strings.Split(val, " ")

	if len(vals) != 2 {
		return "", errors.New("authentication header malformed")
	}

	if vals[0] != "apikey" {
		return "", errors.New("authentication header malformed")
	}

	return vals[1], nil
}
